<template>
  <div class="container">
    <el-row>
      <el-col :span="24">
        <span>タイトル</span>
        <el-input placeholder="Please input" v-model="title"></el-input>
      </el-col>
      <el-col :span="24">
        <span>金額</span>
        <el-radio v-model="radio" label="1">出金</el-radio>
        <el-radio v-model="radio" label="2"> 入金</el-radio>
        <el-input placeholder="Please input" v-model="amount" ></el-input>
      </el-col>
      <el-col :span="24">
        <span>ジャンル</span>
        <el-select style="width: 100%;" v-model="genre" placeholder="please select your zone">
          <el-option label="食費" value="0"></el-option>
          <el-option label="光熱費" value="1"></el-option>
          <el-option label="教育費" value="2"></el-option>
          <el-option label="移動費" value="3"></el-option>
          <el-option label="娯楽費" value="4"></el-option>
          <el-option label="その他" value="5"></el-option>
        </el-select>
      </el-col>
    </el-row>
    <el-row>
      <el-button style="width: 100px;" @click="postData">作成</el-button>
      <el-button style="width: 100px;">
        <nuxt-link to="/">戻る</nuxt-link>
      </el-button>
    </el-row>
  </div>  
</template>

<script>
import axios from "axios"

export default {
  data() {
    return { 
      title: "",
      genre: "",
      amount: "",
      radio: "1"
    }
  },
  methods: {
    postData: function() {
      if(this.radio === "2") {
        this.amount = Number(-this.amount)
        console.log(this.amount)
      }
      axios.post("/api/items",
        { title: this.title, genre: Number(this.genre), amount: Number(this.amount) },
        { headers: { "Content-Type": "application/json" } }
      )
      .then(() => {
        this.$router.push("/")
      })
    }
  }
}
</script>

<style scoped>
.container {
  display: block;
  width: 90%;
  max-width: 800px;
  margin: 0 auto;
  min-height: 100vh;
}

.el-col {
  width: 100%;
  text-align: left;
}

.el-radio {
  margin: 10px 10px;
}

span {
  display: block;
  width: 100%;
  margin-bottom: 10px;
}
</style>