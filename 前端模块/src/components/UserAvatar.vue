<template>
  <div class="avatar-component">
    <div class="avatar-display">
      <img v-if="avatarUrl" :src="avatarUrl" class="avatar-image" alt="用户头像" />
      <div v-else class="avatar-placeholder">
        <span>{{ userInitials }}</span>
      </div>
    </div>
    
    <div class="avatar-actions">
      <el-button type="primary" size="small" @click="showUploadDialog">更换头像</el-button>
    </div>
    
    <!-- 头像上传对话框 -->
    <el-dialog title="更换头像" :visible.sync="dialogVisible" width="400px">
      <div class="upload-container">
        <div class="avatar-preview">
          <img v-if="previewUrl" :src="previewUrl" class="preview-image" />
          <div v-else class="preview-placeholder">
            <span>请选择图片</span>
          </div>
        </div>
        
        <div class="upload-actions">
          <input type="file" ref="fileInput" @change="handleFileChange" accept="image/*" style="display: none" />
          <el-button type="primary" @click="triggerFileInput" :disabled="uploading">选择图片</el-button>
          <el-button type="success" @click="uploadAvatar" :loading="uploading" :disabled="!selectedFile">
            上传头像
          </el-button>
        </div>
        
        <div class="upload-tips">
          <p>支持JPG、PNG格式，文件不超过2MB</p>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import axios from 'axios'
import { mapState } from 'vuex'

export default {
  name: 'UserAvatar',
  data() {
    return {
      dialogVisible: false,
      selectedFile: null,
      previewUrl: '',
      uploading: false
    }
  },
  computed: {
    ...mapState(['user']),
    
    // 获取用户头像URL
    avatarUrl() {
      return this.user?.avatarUrl || ''
    },
    
    // 获取用户名首字母（无头像时显示）
    userInitials() {
      if (!this.user || !this.user.username) return '?'
      return this.user.username.charAt(0).toUpperCase()
    }
  },
  methods: {
    // 显示上传对话框
    showUploadDialog() {
      this.dialogVisible = true
      this.previewUrl = this.avatarUrl
      this.selectedFile = null
    },
    
    // 触发文件选择
    triggerFileInput() {
      this.$refs.fileInput.click()
    },
    
    // 处理文件选择
    handleFileChange(event) {
      const files = event.target.files
      if (files.length > 0) {
        const file = files[0]
        
        // 检查文件类型
        if (!file.type.startsWith('image/')) {
          this.$message.error('请选择图片文件')
          return
        }
        
        // 检查文件大小（2MB限制）
        if (file.size > 2 * 1024 * 1024) {
          this.$message.error('图片大小不能超过2MB')
          return
        }
        
        this.selectedFile = file
        
        // 创建预览
        const reader = new FileReader()
        reader.onload = (e) => {
          this.previewUrl = e.target.result
        }
        reader.readAsDataURL(file)
      }
    },
    
    // 上传头像
    async uploadAvatar() {
      if (!this.selectedFile) {
        this.$message.warning('请先选择图片')
        return
      }
      
      const formData = new FormData()
      formData.append('avatar', this.selectedFile)
      
      this.uploading = true
      try {
        const response = await axios.post('/api/user/avatar', formData, {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        })
        
        // 更新Vuex中的用户信息
        this.$store.commit('setUser', {
          ...this.user,
          avatarUrl: response.data.avatarUrl
        })
        
        this.$message.success('头像上传成功')
        this.dialogVisible = false
      } catch (error) {
        this.$message.error(`上传失败: ${error.response?.data?.error || '未知错误'}`)
        console.error(error)
      } finally {
        this.uploading = false
      }
    }
  }
}
</script>

<style scoped>
.avatar-component {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.avatar-display {
  margin-bottom: 10px;
}

.avatar-image {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid #dcdfe6;
}

.avatar-placeholder {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  background-color: #409eff;
  color: white;
  display: flex;
  justify-content: center;
  align-items: center;
  font-size: 36px;
  font-weight: bold;
}

.upload-container {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.avatar-preview {
  margin-bottom: 20px;
}

.preview-image {
  width: 150px;
  height: 150px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid #dcdfe6;
}

.preview-placeholder {
  width: 150px;
  height: 150px;
  border-radius: 50%;
  background-color: #f5f7fa;
  border: 2px dashed #dcdfe6;
  color: #909399;
  display: flex;
  justify-content: center;
  align-items: center;
}

.upload-actions {
  display: flex;
  gap: 10px;
  margin-bottom: 15px;
}

.upload-tips {
  color: #909399;
  font-size: 12px;
}
</style>
