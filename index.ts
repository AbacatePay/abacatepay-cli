const { createRequire } = require('node:module');
require = createRequire(__filename); 

const yargs = require('yargs');
const localtunnel = require('localtunnel');
const express = require('express');
const { createProxyMiddleware } = require('http-proxy-middleware');

const Express = express.Express;

function getLoggerPrefix(useAbacate: boolean = true, useDate: boolean = true) {
    if(!useAbacate && !useDate) return "";

    return `[${useAbacate && "🥑Abacate CLI"} ${useDate && useAbacate && " - "} ${useDate && new Date().toLocaleString()}]`;
}

function getLogger(useLogger: boolean = true, useAbacate: boolean = true, useDate: boolean = true) {
    if(!useLogger) () => {};

    return (message: string) => {
        console.log(`${getLoggerPrefix(useAbacate, useDate)} ${message}`);
    }
}

(async () => {
    const argv = await yargs
        .option('target', {
            alias: 't',
            type: 'string',
            description: 'URL do servidor local para encaminhamento.',
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
        .help()
        .argv;

    const app = express();
    
    const PORT = 8954;
    const LOCAL_SERVER = argv.target;
    if(!LOCAL_SERVER || LOCAL_SERVER == undefined) {
        console.log('URL do servidor local não informado. Digite abacate --help para ver a lista de opções.');
        process.exit(1);
    }
    const USE_LOGGER = argv.logger;
    const USE_LOG_PREFIX = argv.log_prefix;
    const USE_LOG_TIME = argv.log_time;
    
    const log = getLogger(USE_LOGGER, USE_LOG_PREFIX, USE_LOG_TIME);

    app.use((req: Request, res: Response, next: Function) => {
        log(`🔗 Nova requisição: ${req.method} ${req.url}`);
        next();
    });
    
    app.use('/', (req: Request, res: any, next: Function) => {
        fetch(`${LOCAL_SERVER}${req.url}`, {
            method: req.method,
            headers: req.headers,
            body: req.body
        })
        res.status(200).end('')
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
