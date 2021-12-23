<template>  
  <div>
    <el-container class="title">
      <el-header class="homeHeader">
        <div>&nbsp服务治理平台</div>
        <div>
          <el-row class="demo-avatar demo-basic">
            <el-col :span="12">
              <div class="sub-title"></div>
              <div class="demo-basic--circle">
                <div class="block">
                  <el-avatar :size="40" :src="circleUrl" fit="fill"></el-avatar>
                </div>
              </div>
            </el-col>
          </el-row>
        </div>
      </el-header>
      <el-container class="d">
        <el-aside width="200px">
          <el-menu router class="d">
            <el-submenu
              index="1"
              v-for="(item, index) in this.$router.options.routes"
              :key="index"
              v-if="!item.hidden"
            >
              <template slot="title">
                <i class="el-icon-location"></i>
                <span>{{ item.name }}</span>
              </template>
              <el-menu-item
                :index="children.path"
                v-for="(children, index) in item.children"
                >{{ children.name }}</el-menu-item
              >
            </el-submenu>
          </el-menu>
        </el-aside>
        <el-container>
          <el-main>
            <router-view />
          </el-main>
        </el-container>
      </el-container>
    </el-container>
  </div>
</template>
<script>
export default {
  name: "Home",
  data() {
    return {
      circleUrl:
        "https://gimg2.baidu.com/image_search/src=http%3A%2F%2Fb-ssl.duitang.com%2Fuploads%2Fitem%2F201509%2F20%2F20150920114230_xfMvz.thumb.224_0.gif&refer=http%3A%2F%2Fb-ssl.duitang.com&app=2002&size=f9999,10000&q=a80&n=0&g=0n&fmt=jpeg?sec=1642748405&t=94f948b43fc8071b0c4ba63ee4f643d8",
      squareUrl:
        "https://cube.elemecdn.com/9/c2/f0ee8a3c7c9638a54940382568c9dpng.png",
      sizeList: ["large", "medium", "small"],
    };
  },
};
</script>
<style scoped>
.homeHeader {
  background: rgb(3, 80, 87);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 15px;
  box-sizing: border-box;
  font-size: 20px;
  color: rgb(245, 245, 245);
}
.d {
  background: white;
}
</style>
