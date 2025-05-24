<template>
  <div class="profile-container">
    <h1>个人信息</h1>
    
    <div class="profile-card">
      <div class="avatar-section">
        <UserAvatar />
      </div>
      
      <div class="user-info">
        <el-form label-width="80px">
          <el-form-item label="用户名">
            <span>{{ user.username }}</span>
          </el-form-item>
          <el-form-item label="邮箱">
            <span>{{ user.email || '未设置' }}</span>
          </el-form-item>
          <el-form-item label="注册时间">
            <span>{{ formatDate(user.createdAt) }}</span>
          </el-form-item>
        </el-form>
      </div>
    </div>
    
    <div class="action-buttons">
      <el-button type="primary" @click="goToFiles">
        <i class="el-icon-folder"></i> 管理我的文件
      </el-button>
      <el-button @click="goToHome">
        <i class="el-icon-back"></i> 返回首页
      </el-button>
    </div>
  </div>
</template>

<script>
import { mapState } from 'vuex'
import UserAvatar from '@/components/UserAvatar.vue'

export default {
  name: 'Profile',
  components: {
    UserAvatar
  },
  computed: {
    ...mapState(['user'])
  },
  methods: {
    formatDate(dateString) {
      if (!dateString) return '未知'
      const date = new Date(dateString)
      return date.toLocaleString()
    },
    goToFiles() {
      this.$router.push('/files')
    },
    goToHome() {
      this.$router.push('/home')
    }
  }
}
</script>

<style scoped>
.profile-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.profile-card {
  display: flex;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  padding: 30px;
  margin-bottom: 30px;
}

.avatar-section {
  margin-right: 40px;
}

.user-info {
  flex: 1;
}

.action-buttons {
  display: flex;
  justify-content: center;
  gap: 20px;
  margin-top: 20px;
}

@media (max-width: 768px) {
  .profile-card {
    flex-direction: column;
  }
  
  .avatar-section {
    margin-right: 0;
    margin-bottom: 20px;
    display: flex;
    justify-content: center;
  }
}
</style>
