<template>
  <div class="file-manager">
    <h1>文件管理</h1>
    
    <!-- 文件上传区域 -->
    <div class="upload-area">
      <h2>上传文件</h2>
      <div class="upload-box" @click="triggerFileInput" @dragover.prevent @drop.prevent="handleFileDrop">
        <input type="file" ref="fileInput" @change="handleFileChange" style="display: none" />
        <div class="upload-icon">
          <i class="el-icon-upload"></i>
        </div>
        <div class="upload-text">
          <span>点击或拖拽文件到此处上传</span>
          <p>支持jpg, png, pdf等常见格式，单个文件不超过10MB</p>
        </div>
      </div>
      <el-button type="primary" @click="uploadFile" :loading="uploading" :disabled="!selectedFile">
        上传文件
      </el-button>
    </div>

    <!-- 文件列表区域 -->
    <div class="file-list">
      <h2>我的文件</h2>
      <el-table :data="fileList" style="width: 100%" v-loading="loading">
        <el-table-column label="文件名" prop="fileName"></el-table-column>
        <el-table-column label="类型" width="120">
          <template slot-scope="scope">
            <span>{{ getFileTypeLabel(scope.row.fileType) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="大小" width="120">
          <template slot-scope="scope">
            <span>{{ formatFileSize(scope.row.fileSize) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="上传时间" width="180" prop="uploadAt"></el-table-column>
        <el-table-column label="操作" width="200">
          <template slot-scope="scope">
            <el-button type="text" @click="previewFile(scope.row)">预览</el-button>
            <el-button type="text" @click="downloadFile(scope.row)">下载</el-button>
            <el-button type="text" class="danger-text" @click="deleteFile(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 文件预览对话框 -->
    <el-dialog title="文件预览" :visible.sync="previewVisible" width="70%">
      <div class="preview-container" v-if="currentFile">
        <div v-if="isImage(currentFile.fileType)">
          <img :src="currentFile.fileUrl" class="preview-image" />
        </div>
        <div v-else-if="isPdf(currentFile.fileType)">
          <iframe :src="currentFile.fileUrl" class="preview-frame"></iframe>
        </div>
        <div v-else class="preview-unsupported">
          <p>此文件类型不支持预览，请下载后查看</p>
          <el-button type="primary" @click="downloadFile(currentFile)">下载文件</el-button>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'FileManager',
  data() {
    return {
      fileList: [],
      selectedFile: null,
      uploading: false,
      loading: false,
      previewVisible: false,
      currentFile: null
    }
  },
  created() {
    this.fetchFileList()
  },
  methods: {
    // 获取文件列表
    async fetchFileList() {
      this.loading = true
      try {
        const response = await axios.get('/api/files')
        this.fileList = response.data.files || []
      } catch (error) {
        this.$message.error('获取文件列表失败')
        console.error(error)
      } finally {
        this.loading = false
      }
    },
    
    // 触发文件选择
    triggerFileInput() {
      this.$refs.fileInput.click()
    },
    
    // 处理文件选择
    handleFileChange(event) {
      const files = event.target.files
      if (files.length > 0) {
        this.selectedFile = files[0]
        this.$message.success(`已选择文件: ${this.selectedFile.name}`)
      }
    },
    
    // 处理文件拖放
    handleFileDrop(event) {
      const files = event.dataTransfer.files
      if (files.length > 0) {
        this.selectedFile = files[0]
        this.$message.success(`已选择文件: ${this.selectedFile.name}`)
      }
    },
    
    // 上传文件
    async uploadFile() {
      if (!this.selectedFile) {
        this.$message.warning('请先选择文件')
        return
      }
      
      // 检查文件大小
      if (this.selectedFile.size > 10 * 1024 * 1024) {
        this.$message.error('文件大小不能超过10MB')
        return
      }
      
      const formData = new FormData()
      formData.append('file', this.selectedFile)
      
      this.uploading = true
      try {
        const response = await axios.post('/api/file/upload', formData, {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        })
        
        this.$message.success('文件上传成功')
        this.selectedFile = null
        this.$refs.fileInput.value = ''
        this.fetchFileList() // 刷新文件列表
      } catch (error) {
        this.$message.error(`上传失败: ${error.response?.data?.error || '未知错误'}`)
        console.error(error)
      } finally {
        this.uploading = false
      }
    },
    
    // 预览文件
    previewFile(file) {
      this.currentFile = file
      this.previewVisible = true
    },
    
    // 下载文件
    downloadFile(file) {
      window.open(file.fileUrl, '_blank')
    },
    
    // 删除文件
    async deleteFile(file) {
      this.$confirm('确定要删除此文件吗?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async () => {
        try {
          const fileName = file.fileUrl.split('/').pop()
          await axios.post(`/api/file/delete/${fileName}`)
          this.$message.success('文件删除成功')
          this.fetchFileList() // 刷新文件列表
        } catch (error) {
          this.$message.error(`删除失败: ${error.response?.data?.error || '未知错误'}`)
          console.error(error)
        }
      }).catch(() => {
        // 取消删除
      })
    },
    
    // 格式化文件大小
    formatFileSize(size) {
      if (size < 1024) {
        return size + ' B'
      } else if (size < 1024 * 1024) {
        return (size / 1024).toFixed(2) + ' KB'
      } else {
        return (size / (1024 * 1024)).toFixed(2) + ' MB'
      }
    },
    
    // 获取文件类型标签
    getFileTypeLabel(type) {
      if (type.startsWith('image/')) {
        return '图片'
      } else if (type === 'application/pdf') {
        return 'PDF'
      } else if (type.includes('word')) {
        return 'Word'
      } else if (type.includes('excel')) {
        return 'Excel'
      } else if (type === 'text/plain') {
        return '文本'
      } else {
        return '其他'
      }
    },
    
    // 判断是否为图片
    isImage(type) {
      return type.startsWith('image/')
    },
    
    // 判断是否为PDF
    isPdf(type) {
      return type === 'application/pdf'
    }
  }
}
</script>

<style scoped>
.file-manager {
  padding: 20px;
}

.upload-area {
  margin-bottom: 30px;
  padding: 20px;
  background-color: #f5f7fa;
  border-radius: 4px;
}

.upload-box {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 200px;
  border: 2px dashed #dcdfe6;
  border-radius: 6px;
  cursor: pointer;
  margin-bottom: 20px;
  transition: all 0.3s;
}

.upload-box:hover {
  border-color: #409eff;
}

.upload-icon {
  font-size: 48px;
  color: #909399;
  margin-bottom: 10px;
}

.upload-text {
  text-align: center;
}

.upload-text span {
  font-size: 16px;
  color: #606266;
}

.upload-text p {
  margin-top: 10px;
  color: #909399;
  font-size: 14px;
}

.preview-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 300px;
}

.preview-image {
  max-width: 100%;
  max-height: 500px;
}

.preview-frame {
  width: 100%;
  height: 500px;
  border: none;
}

.preview-unsupported {
  text-align: center;
  padding: 30px;
}

.danger-text {
  color: #f56c6c;
}
</style>
