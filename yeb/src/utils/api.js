import axios from "axios";
import { Message } from 'element-ui';
import router from '../router'

// 响应拦截器【统一 fang'b】
axios.interceptors.response.use(success=>{
    // 业务逻辑错误
    if(success.status && success.status==200){
        if(success.data.code==500 || success.data.code==401 || success.data.code==403){
            Message.error({message:success.data.message})
            return
        }
        if(success.data.message){
            Message.success({message:success.data.message})
        }
    }
    return success.data;
},error=>{
    // 连接口都没有访问到
    if(error.response.code==504||error.response.code==404){
        Message.error({message:'服务器被吃了'})
    }else if(error.response.status ==403){
        Message.error({message:'权限不足'})
    }else if(error.response.status ==401){
        Message.error({message:'尚未登录'})
        router.replace('/')
    }else {
        if(error.response.data){
            Message.error({message:error.response.data.message})
        }else{
            Message.error({message:'位置错误！'})
        }
    }
    return;
});

let base = ''
// 传送json格式的post请求
export const postRequest = (url,params)=>{
    return axios({
        method:'post',
        url:`${base}${url}`,
        data: params
    })
}