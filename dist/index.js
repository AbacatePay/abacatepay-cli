"use strict";
const yargs = require('yargs');
const localtunnel = require('localtunnel');
const express = require('express');
const { createProxyMiddleware } = require('http-proxy-middleware');
const Express = express.Express;
function getLoggerPrefix(useAbacate = true, useDate = true) {
    if (!useAbacate && !useDate)
        return "";
    return `[${useAbacate && "🥑Abacate CLI"} ${useDate && useAbacate && " - "} ${useDate && new Date().toLocaleString()}]`;
}
function getLogger(useLogger = true, useAbacate = true, useDate = true) {
    if (!useLogger)
        () => { };
    return (message) => {
        console.log(`${getLoggerPrefix(useAbacate, useDate)} ${message}`);
    };
}
(async () => {
    const argv = await yargs
        .option('target', {
        alias: 't',
        type: 'string',
        description: 'URL do servidor local para encaminhamento.',
        requiresArg: true
    })
        .option('logger', {
        alias: 'l',
        type: 'boolean',
        description: 'Usar o logger.',
        default: true
    })
        .option('log_prefix', {
        alias: 'lp',
        type: 'boolean',
        description: 'Usar o prefixo "🥑Abacate CLI" no logger.',
        default: true
    })
        .option('log_time', {
        alias: 'lt',
        type: 'boolean',
        description: 'Usar o tempo no prefixo do logger.',
        default: true
    })
        .option('logger', {
        alias: 'l',
        type: 'boolean',
        description: 'Usar o logger',
        default: true
    })
        .help()
        .argv;
    const app = express();
    const PORT = 8954;
    const LOCAL_SERVER = argv.target;
    const USE_LOGGER = argv.logger;
    const USE_LOG_PREFIX = argv.log_prefix;
    const USE_LOG_TIME = argv.log_time;
    const log = getLogger(USE_LOGGER, USE_LOG_PREFIX, USE_LOG_TIME);
    app.use((req, res, next) => {
        log(`🔗 Nova requisição: ${req.method} ${req.url}`);
        next();
    });
    app.use('/', (req, res, next) => {
        fetch(`${LOCAL_SERVER}${req.url}`, {
            method: req.method,
            headers: req.headers,
            body: req.body
        });
        res.status(200).end('');
    });
    app.listen(PORT, async () => {
        log(`🚀 Servidor rodando na porta ${PORT}`);
        log(`🔄 Encaminhando requisições para ${LOCAL_SERVER}`);
        const tunnel = await localtunnel({ port: PORT, subdomain: '', allow_invalid_cert: true });
        log(`🌍 Servidor acessível publicamente em: ${tunnel.url}`);
        log(`🔑 Coloque a URL do servidor público em sua https://www.abacatepay.com`);
        tunnel.on('close', () => {
            log('❌ Túnel fechado.');
        });
    });
})();
