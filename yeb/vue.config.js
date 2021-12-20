let proxyObj = {}

proxyObj['/']={
    //websocket
    ws:false,
    target:"http://localhost:8080",
    //发送请求头的host会被设置为target
    changeOrigin: true,
    //不重写请求地址
    pathReWrite:{
        '^/':'/'
    }
}

module.exports={
    devServer:{
        host:'localhost',
        port:8082,
        proxy: proxyObj
    }
}