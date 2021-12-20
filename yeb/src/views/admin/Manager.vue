<template>
    <div>
        <div>
            <el-select v-model="value" placeholder="请选择工作空间" @change="getServiceStatus">
            <el-option
            v-for="item in options"
            :key="item.value"
            :label="item.label"
            :value="item.value">
            </el-option>
            </el-select>
        </div>
        <div id="ServiceStatus" style="width: 400px;height: 300px;"></div>
        <!-- <div id="NodeStatus" style="width: 400px;height: 300px;"></div> -->
    </div>
</template>

<script>
import Axios from 'axios';
export default {
  data() {
    return {
        options: [{
          value: 'default',
          label: 'default'
        }, {
          value: '选项2',
          label: '1111'
        }, {
          value: '选项3',
          label: '2222'
        }, {
          value: '选项4',
          label: '3333'
        }, {
          value: '选项5',
          label: '4444'
        }],
        value: '',
      charts: "",
      ServiceStatusOpinion: ["服务健康","部分不可用","不可用"],
      ServiceStatusOpinionData: []
    };
  },
  mounted() {
    this.$nextTick(function() {
        // this.getServiceStatus("ServiceStatus",this.value);
        //  this.getServiceStatus("NodeStatus","default");
        // this.drawPie("ServiceStatus");
        // this.drawPie("NodeStatus");
    });
  },
  methods: {
    clearData(){
        this.ServiceStatusOpinionData = []
    },
    getServiceStatus(value){
        this.clearData()
        // var api = 'http://202.38.247.217:8080/api/v1/rest/svc/workspace/'+namespace+'/status';
        var api = 'http://localhost:8080/api/v1/rest/svc/workspace/'+value+'/status';
        Axios.get(api).then(res=>{
            console.log(res.data.data)
            this.ServiceStatusOpinionData.push( { value: res.data.data["healthy"], name: "服务健康"})
            this.ServiceStatusOpinionData.push( { value: res.data.data["unhealthy"], name: "部分不可用"})
            this.ServiceStatusOpinionData.push( { value: res.data.data["error"], name: "不可用"})
            this.charts = this.$echarts.init(document.getElementById("ServiceStatus"));
      this.charts.setOption({
        tooltip: {
          trigger: "item",
          formatter: "{a}<br/>{b}:{c} ({d}%)"
        },
        legend: {
          bottom: 10,
          left: "center",
          data: this.ServiceStatusOpinion
        },
        series: [
          {
            name: "状态",
            type: "pie",
            radius: "65%",
            center: ["50%", "50%"],
            avoidLabelOverlap: false,
            itemStyle: {
              emphasis: {
                shadowBlur: 10,
                shadowOffsetX: 0,
                shadowColor: "rgba(0, 0, 0, 0.5)"
              },
              color: function(params) {
                //自定义颜色
                var colorList = ["#1ab394", "#79d2c0","#e9321e"];
                return colorList[params.dataIndex];
              }
            },
            data: this.ServiceStatusOpinionData
          }
        ]
      });
        }).catch(error=>{
            // console.log('Error',error.message);
        })
    },
    getNodeStatus(id){

    }
  },
};
</script>

<style  scoped>
    .pie-wrap {
        width: 50%;
        height: 126px;
    }
</style>

