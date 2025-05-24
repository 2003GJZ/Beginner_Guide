<template>
  <div class="login-container">
    <el-card class="login-card">
      <div class="title">
        <h2>{{ isLogin ? '用户登录' : '用户注册' }}</h2>
      </div>
      
      <el-form :model="formData" :rules="rules" ref="loginForm" label-width="80px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="formData.username" placeholder="请输入用户名"></el-input>
        </el-form-item>
        
        <el-form-item label="密码" prop="password">
          <el-input v-model="formData.password" type="password" placeholder="请输入密码"></el-input>
        </el-form-item>
        
        <el-form-item v-if="!isLogin" label="确认密码" prop="confirmPassword">
          <el-input v-model="formData.confirmPassword" type="password" placeholder="请再次输入密码"></el-input>
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="submitForm" :loading="loading">{{ isLogin ? '登录' : '注册' }}</el-button>
          <el-button @click="switchMode">{{ isLogin ? '去注册' : '去登录' }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script>
export default {
  name: 'Login',
  data() {
    // 自定义校验规则：确认密码
    const validateConfirmPassword = (rule, value, callback) => {
      if (value !== this.formData.password) {
        callback(new Error('两次输入的密码不一致'))
      } else {
        callback()
      }
    }
    
    return {
      // 是否为登录模式
      isLogin: true,
      // 加载状态
      loading: false,
      // 表单数据
      formData: {
        username: '',
        password: '',
        confirmPassword: ''
      },
      // 表单验证规则
      rules: {
        username: [
          { required: true, message: '请输入用户名', trigger: 'blur' },
          { min: 3, max: 20, message: '用户名长度应为3-20个字符', trigger: 'blur' }
        ],
        password: [
          { required: true, message: '请输入密码', trigger: 'blur' },
          { min: 6, max: 20, message: '密码长度应为6-20个字符', trigger: 'blur' }
        ],
        confirmPassword: [
          { required: true, message: '请再次输入密码', trigger: 'blur' },
          { validator: validateConfirmPassword, trigger: 'blur' }
        ]
      }
    }
  },
  methods: {
    // 切换登录/注册模式
    switchMode() {
      this.isLogin = !this.isLogin
      this.$refs.loginForm.resetFields()
    },
    
    // 提交表单
    submitForm() {
      this.$refs.loginForm.validate(async valid => {
        if (!valid) return
        
        this.loading = true
        
        try {
          if (this.isLogin) {
            // 登录操作
            await this.$store.dispatch('login', {
              username: this.formData.username,
              password: this.formData.password
            })
            
            // 登录成功，跳转到首页
            this.$router.push('/home')
            this.$message.success('登录成功')
          } else {
            // 注册操作
            console.log('准备提交注册数据:', {
              username: this.formData.username,
              password: this.formData.password
            });
            
            // 确保密码字段正确传递
            const registerData = {
              username: this.formData.username,
              password: this.formData.password
            };
            
            await this.$store.dispatch('register', registerData)
            
            // 注册成功，切换到登录模式
            this.isLogin = true
            this.$refs.loginForm.resetFields()
            this.$message.success('注册成功，请登录')
          }
        } catch (error) {
          // 处理错误
          const errorMsg = error.response?.data?.error || (this.isLogin ? '登录失败' : '注册失败')
          this.$message.error(errorMsg)
        } finally {
          this.loading = false
        }
      })
    }
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: #f5f7fa;
}

.login-card {
  width: 400px;
  padding: 20px;
}

.title {
  text-align: center;
  margin-bottom: 20px;
}
</style>
