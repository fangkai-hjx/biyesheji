<template>
  <el-dialog
    :title="svc + '  服务详情'"
    :close-on-click-modal="false"
    :visible.sync="visible"
    center="center"
    width="1200px"
    @closed="dialogClose"
  >
    <div
      class="float-div"
      id="ServiceDetail"
      style="width: 300px; height: 300px"
    ></div>
    <el-descriptions
      class="float-div"
      title="基础信息"
      direction="vertical"
      :column="5"
      border
      size="small"
    >
      <el-descriptions-item label="工作空间名称">{{
        svc
      }}</el-descriptions-item>
      <el-descriptions-item label="服务名称">{{ svc }}</el-descriptions-item>
      <el-descriptions-item label="服务端口">
        <el-tag size="small">{{ targetPort }}</el-tag>
      </el-descriptions-item>
      <el-descriptions-item label="容器端口">
        <el-tag size="small">{{ containerPort }}</el-tag>
      </el-descriptions-item>
      <el-descriptions-item label="协议">{{ prot }}</el-descriptions-item>
    </el-descriptions>
    <el-descriptions
      class="float-div"
      title="元信息"
      direction="vertical"
      :column="5"
      border
      size="small"
    >
      <el-descriptions-item label="服务类型">{{ Type }}</el-descriptions-item>
      <el-descriptions-item label="服务 IP">
        {{ ClusterIP }}
      </el-descriptions-item>
      <el-descriptions-item label="标签选择器">
        <el-tag size="small">{{ Selector }}</el-tag>
      </el-descriptions-item>
      <el-descriptions-item label="注解">
        <el-tag size="small">ttt</el-tag>
      </el-descriptions-item>
    </el-descriptions>
    <el-table :data="tableData" border width="100%">
      <el-table-column
        prop="name"
        label="实例名"
        width="200px"
      ></el-table-column>
      <el-table-column
        prop="pod_ip"
        label="实例IP"
        width="130px"
      ></el-table-column>
      <el-table-column
        prop="host_ip"
        label="主机IP"
        width="130px"
      ></el-table-column>
      <el-table-column
        prop="image"
        label="镜像"
        width="200px"
      ></el-table-column>
      <el-table-column label="状态" width="100px">
        <el-table-column prop="PodScheduled" label="调度">
          <template slot-scope="scope">
            <el-switch
              v-model="scope.row.PodScheduled"
              active-color="#13ce66"
              inactive-color="#ff4949"
              disabled
            >
            </el-switch>
          </template>
        </el-table-column>
        <el-table-column prop="Initialized" label="初始化">
          <template slot-scope="scope">
            <el-switch
              v-model="scope.row.Initialized"
              active-color="#13ce66"
              inactive-color="#ff4949"
              disabled
            >
            </el-switch>
          </template>
        </el-table-column>
        <el-table-column prop="ContainersReady" label="容器启动">
          <template slot-scope="scope">
            <el-switch
              v-model="scope.row.ContainersReady"
              active-color="#13ce66"
              inactive-color="#ff4949"
              disabled
            >
            </el-switch>
          </template>
        </el-table-column>
        <el-table-column prop="Ready" label="Ready">
          <template slot-scope="scope">
            <el-switch
              v-model="scope.row.Ready"
              active-color="#13ce66"
              inactive-color="#ff4949"
              disabled
            >
            </el-switch>
          </template>
        </el-table-column>
      </el-table-column>
      <el-table-column
        fixed="right"
        header-align="center"
        align="center"
        label="实例详情"
      >
        <!-- <template slot-scope="scope">
          <el-button
            type="button"
            size="small"
            @click="getServiceDetails(scope.row.name)"
            >查看</el-button
          >
        </template> -->
      </el-table-column>
    </el-table>
  </el-dialog>
</template>

<script>
import Axios from "axios";
import global from '../common.vue'
export default {
  data() {
    return {
      tableData: [],
      visible: false,
      svc: "",
      ns: "",

      ServiceDetailOpinion: [],
      ServiceDetailOpinionData: [],

      prot: "",
      targetPort: "",
      containerPort: "",
      ClusterIP: "",
      Selector: "Unknow",
      Type: "Unknow",
      Annotations: "Unknow",
    };
  },

  methods: {
    getPodChart() {},
    dialogClose() {
      this.catelogPath = [];
    },
    clear() {
      this.tableData = [];
      this.ServiceDetailOpinion = [];
      this.ServiceDetailOpinionData = [];
    },
    getCategorys() {},
    init(servename, namespace) {
      console.log("init");
      this.svc = servename;
      this.ns = namespace;
      this.visible = true;
      // console.log(this.svc,this.ns)
      // console.log("getServiceDetail")
      this.getServiceDetail();
    },
    getServiceDetail() {
      this.clear();
      var api =
        global.httpUrl+"/api/v1/rest/svc/workspace/" +
        this.ns +
        "/" +
        this.svc +
        "/";
      // 根据下拉框加载工作空间的数据
      Axios.get(api)
        .then((res) => {
          var p = "pods";
          if (res.data.data["success"] != 0) {
            this.ServiceDetailOpinion.push("success");
            this.ServiceDetailOpinionData.push({
              value: res.data.data["success"],
              name: "success",
            });
          }
          if (res.data.data["fail"] != 0) {
            this.ServiceDetailOpinion.push("fail");
            this.ServiceDetailOpinionData.push({
              value: res.data.data["fail"],
              name: "fail",
            });
          }
          this.prot = res.data.data["protocol"];
          this.targetPort = res.data.data["port"];
          this.containerPort = res.data.data["container_port"];
          this.ClusterIP = res.data.data["cluster_ip"];
          if (res.data.data["selector"] != "") {
            this.Selector = res.data.data["selector"];
          }
          if (res.data.data["type"] != "") {
            this.Type = res.data.data["type"];
          }
          if (res.data.data["annotations"] != null) {
            this.Annotations = res.data.data["annotations"];
          }
          for (const key in res.data.data[p]) {
            this.tableData.push({
              name: res.data.data[p][key]["name"],
              Initialized: res.data.data[p][key]["initialized"],
              ContainersReady: res.data.data[p][key]["containersReady"],
              Ready: res.data.data[p][key]["ready"],
              PodScheduled: res.data.data[p][key]["podScheduled"],
              image: res.data.data[p][key]["image"],
              host_ip: res.data.data[p][key]["host_ip"],
              pod_ip: res.data.data[p][key]["pod_ip"],
            });
          }
          this.charts = this.$echarts.init(
            document.getElementById("ServiceDetail")
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
              data: this.ServiceDetailOpinion,
            },
            series: [
              {
                name: "状态",
                type: "pie",
                radius: "50%",
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
                    var colorList = ["#1ab394", "#ff2d51", "#e9321e"];
                    return colorList[params.dataIndex];
                  },
                },
                data: this.ServiceDetailOpinionData,
              },
            ],
          });
        })
        .catch((error) => {
          console.log("Error", error.message);
        });
    },
  },
  created() {},
};
</script>
<style scoped>
.float-div {
  float: left;
  margin: 50px;
}
</style>
