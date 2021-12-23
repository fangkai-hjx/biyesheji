<template>
  <div class="mod-config">
    <el-form
      :inline="true"
      :model="dataForm"
      @keyup.enter.native="getDataList()"
    >
      <el-form-item>
        <el-select
          v-model="value"
          placeholder="请选择工作空间"
          @change="getServiceStatus"
        >
          <el-option
            v-for="item in options"
            :key="item.value"
            :label="item.label"
            :value="item.value"
          >
          </el-option>
        </el-select>
        工作空间
      </el-form-item>
      <el-form-item>
        <el-input
          v-model="dataForm.key"
          placeholder="请输入服务名称"
          clearable
        ></el-input>
      </el-form-item>
      <el-form-item>
        <el-button @click="getDataList()">查询</el-button>
        <!-- <el-button v-if="isAuth('member:member:save')" type="primary" @click="addOrUpdateHandle()">新增</el-button>
        <el-button v-if="isAuth('member:member:delete')" type="danger" @click="deleteHandle()" :disabled="dataListSelections.length <= 0">批量删除</el-button>-->
      </el-form-item>
      <div style="float:left">
        <div
          id="ServiceStatus"
          style="width: 600px; height: 300px"
          class="float-div"
        ></div>
        <div
          id="PodStatus"
          style="width: 600px; height: 300px"
          class="float-div"
        ></div>
      </div>
    </el-form>
    <el-table :data="tableData" border width="100%" :cell-style="cellStyle">
      <el-table-column
        prop="name"
        label="服务名"
        width="250px"
      ></el-table-column>
      <el-table-column prop="runtime" label="运行时长(小时)"></el-table-column>
      <el-table-column prop="cluster_ip" label="服务IP"></el-table-column>
      <el-table-column prop="total" label="实例总数"></el-table-column>
      <el-table-column prop="success" label="可用实例"></el-table-column>
      <el-table-column prop="fail" label="故障实例"></el-table-column>
      <el-table-column prop="status" label="服务状态"></el-table-column>
      <el-table-column prop="session_affinity" label="会话亲和性">
        <template slot-scope="scope">
          <el-switch
            v-model="scope.row.session_affinity"
            active-color="#13ce66"
            inactive-color="#ff4949"
          >
          </el-switch>
        </template>
      </el-table-column>
      <el-table-column prop="success_lu" label="可用率"> </el-table-column>
      <el-table-column
        fixed="right"
        header-align="center"
        align="center"
        label="服务详情"
      >
        <template slot-scope="scope">
          <el-button
            type="button"
            size="small"
            @click="getServiceDetails(scope.row.name)"
            >查看</el-button
          >
        </template>
      </el-table-column>
    </el-table>
    <!-- <div id="NodeStatus" style="width: 400px;height: 300px;"></div> -->
    <!-- 弹窗, 新增 / 修改 -->
    <add-or-update v-if="addOrUpdateVisible" ref="addOrUpdate"></add-or-update>
  </div>
</template>

<script>
import Axios from "axios";
import global from '../common.vue'
import AddOrUpdate from "./service-pod-detail.vue";
export default {
  components: { AddOrUpdate },
  data() {
    return {
      options: [],
      value: "",
      charts: "dadad",
      ServiceStatusOpinion: [],
      ServiceStatusOpinionData: [],
      tableData: [],

      PodStatusOpinion: [],
      PodStatusOpinionData: [],
      dataForm: {
        key: "",
      },
      addOrUpdateVisible: false,
    };
  },
  mounted() {
    // 获取下拉框的数据
    this.getSelectData();
    this.$nextTick(function () {
      // this.getServiceStatus("ServiceStatus",this.value);
      //  this.getServiceStatus("NodeStatus","default");
      // this.drawPie("ServiceStatus");
      // this.drawPie("NodeStatus");
    });
  },
  methods: {
    cellStyle(row, column, rowIndex, columnIndex) {
      if (row.column.label === "服务状态" && row.row.status === "ERROR") {
        return "background:#CD9B9B;color:whitesmoke;font-size:15px;";
      } else if (
        row.column.label === "服务状态" &&
        row.row.status === "UNHEALTHY"
      ) {
        return "background:#DAA520;color:whitesmoke;font-size:15px;";
      } else if (
        row.column.label === "服务状态" &&
        row.row.status === "HEALTHY"
      ) {
        return "background:	#548B54;color:whitesmoke;font-size:15px;";
      }
    },
    clearData() {
      this.ServiceStatusOpinionData = [];
      this.tableData = [];

      this.PodStatusOpinion = [];
      this.PodStatusOpinionData = [];
    },

    // 下拉选项 加载 命名空间
    getSelectData() {
      var api = global.httpUrl+"/api/v1/rest/ns/workspace/all";
      // 根据下拉框加载工作空间的数据
      Axios.get(api)
        .then((res) => {
          for (const key in res.data.data) {
            this.options.push({
              value: res.data.data[key]["name"],
              label: res.data.data[key]["name"],
            });
          }
        })
        .catch((error) => {
          // console.log('Error',error.message);
        });
    },
    getPodChart() {
      var api =
        global.httpUrl+"/api/v1/rest/svc/workspace/" +
        this.value +
        "/pod/status";
      Axios.get(api)
        .then((res) => {
          console.log(res.data.data);
          for (const key in res.data.data) {
            if (res.data.data[key] != 0) {
              this.PodStatusOpinion.push(key);
              this.PodStatusOpinionData.push({
                value: res.data.data[key],
                name: key,
              });
            }
          }
          this.charts = this.$echarts.init(
            document.getElementById("PodStatus")
          );
          this.charts.setOption({
            title: {
              text: "实例运行状况",
              subtext: "",
              left: "center",
            },
            tooltip: {
              trigger: "item",
              formatter: "{a}<br/>{b}:{c} ({d}%)",
            },
            legend: {
              bottom: 10,
              left: "center",
              data: this.PodStatusOpinion,
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
                    shadowColor: "rgba(0, 0, 0, 0.5)",
                  },
                  color: function (params) {
                    //自定义颜色
                    var colorList = ["#ff2d51", "#1ab394", "#e9321e"];
                    return colorList[params.dataIndex];
                  },
                },
                data: this.PodStatusOpinionData,
              },
            ],
          });
        })
        .catch((error) => {
          // console.log('Error',error.message);
        });
    },
    getServiceDetails(name) {
      this.addOrUpdateVisible = true;
      this.$nextTick(() => {
        this.$refs.addOrUpdate.init(name, this.value);
      });
    },

    getServiceChart() {
      var api =
        global.httpUrl+"/api/v1/rest/svc/workspace/" +
        this.value +
        "/svc/status";
      Axios.get(api)
        .then((res) => {
          console.log(res.data.data);
          for (const key in res.data.data) {
            if (res.data.data[key] != 0) {
              this.ServiceStatusOpinion.push(key);
              this.ServiceStatusOpinionData.push({
                value: res.data.data[key],
                name: key,
              });
            }
          }
          this.charts = this.$echarts.init(
            document.getElementById("ServiceStatus")
          );
          this.charts.setOption({
            title: {
              text: "运行状况监控",
              subtext: "",
              left: "center",
            },
            tooltip: {
              trigger: "item",
              formatter: "{a}<br/>{b}:{c} ({d}%)",
            },
            legend: {
              bottom: 10,
              left: "center",
              data: this.ServiceStatusOpinion,
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
                    shadowColor: "rgba(0, 0, 0, 0.5)",
                  },
                  color: function (params) {
                    //自定义颜色
                    var colorList = ["#1ab394", "#79d2c0", "#e9321e"];
                    return colorList[params.dataIndex];
                  },
                },
                data: this.ServiceStatusOpinionData,
              },
            ],
          });
        })
        .catch((error) => {
          // console.log('Error',error.message);
        });
    },
    getServiceTable() {
      var api =
        global.httpUrl+"/api/v1/rest/svc/workspace/" +
        this.value +
        "/all";
      console.log("getServiceTable", this.value);
      Axios.get(api)
        .then((res) => {
          for (const key in res.data.data) {
            this.tableData.push({
              name: res.data.data[key]["name"],
              cluster_ip: res.data.data[key]["cluster_ip"],
              session_affinity: res.data.data[key]["session_affinity"],
              status: res.data.data[key]["status"],
              runtime: res.data.data[key]["runtime"],
              success: res.data.data[key]["success"],
              fail: res.data.data[key]["fail"],
              total: res.data.data[key]["total"],
              success_lu: res.data.data[key]["success_lu"] + " %",
            });
          }
        })
        .catch((error) => {
          // console.log('Error',error.message);
        });
    },
    getServiceStatus() {
      this.clearData();
      this.getServiceTable();
      this.getData();
    },
    getData() {
      this.getServiceChart();
      this.getPodChart();
    },
    getNodeStatus(id) {},
  },
};
</script>

<style  scoped>
/* .bottom {
  margin-top: 13px;
  line-height: 12px;
} */
.float-div {
  float: left;
  margin: 30px;
}
/deep/ .el-table .warning-row {
  background: #f7edef;
}

/deep/ .el-table .success-row {
  background: white;
}
/deep/ .el-table .error-row {
  background: rgb(241, 210, 210);
}
</style>

