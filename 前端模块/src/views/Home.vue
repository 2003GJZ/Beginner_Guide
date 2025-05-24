<template>
  <div class="home-container">
    <el-container>
      <!-- 头部 -->
      <el-header>
        <div class="header-content">
          <div class="logo-section">
            <h2>任务管理系统</h2>
          </div>
          <div class="nav-links">
            <el-button type="text" @click="$router.push('/home')">
              <i class="el-icon-s-home"></i> 首页
            </el-button>
            <el-button type="text" @click="$router.push('/files')">
              <i class="el-icon-folder"></i> 文件管理
            </el-button>
          </div>
          <div class="user-info">
            <div class="username-container" @click="$router.push('/profile')">
              <div v-if="user && user.avatarUrl" class="avatar-container">
                <img :src="user.avatarUrl" class="mini-avatar" alt="头像">
              </div>
              <div v-else class="avatar-placeholder">{{ userInitials }}</div>
              <span class="username-text">{{ user ? user.username : '用户' }}</span>
            </div>
            <el-button type="danger" size="small" @click="logout" class="logout-btn">
              <i class="el-icon-switch-button"></i> 退出
            </el-button>
          </div>
        </div>
      </el-header>
      
      <!-- 主体内容 -->
      <el-main>
        <div class="task-header">
          <h3>我的任务列表</h3>
          <el-button type="primary" size="small" @click="showAddTaskDialog">新建任务</el-button>
        </div>
        
        <!-- 任务过滤器 -->
        <div class="task-filter">
          <div class="filter-row">
            <span class="filter-label">状态：</span>
            <el-radio-group v-model="taskFilter" size="small">
              <el-radio-button label="all">全部</el-radio-button>
              <el-radio-button label="active">未完成</el-radio-button>
              <el-radio-button label="completed">已完成</el-radio-button>
            </el-radio-group>
          </div>
          
          <div class="filter-row">
            <span class="filter-label">优先级：</span>
            <el-radio-group v-model="priorityFilter" size="small">
              <el-radio-button label="all">全部</el-radio-button>
              <el-radio-button label="low">低</el-radio-button>
              <el-radio-button label="medium">中</el-radio-button>
              <el-radio-button label="high">高</el-radio-button>
            </el-radio-group>
          </div>
        </div>
        
        <!-- 任务列表 -->
        <el-table
          v-loading="loading"
          :data="filteredTasks"
          style="width: 100%"
          empty-text="暂无任务"
        >
          <el-table-column prop="title" label="任务名称" min-width="180">
            <template slot-scope="scope">
              <el-checkbox
                v-model="scope.row.completed"
                @change="updateTaskStatus(scope.row)"
              ></el-checkbox>
              <span :class="{ 'task-completed': scope.row.completed }">{{ scope.row.title }}</span>
            </template>
          </el-table-column>
          
          <el-table-column prop="description" label="任务描述" min-width="200"></el-table-column>
          
          <el-table-column prop="priority" label="优先级" width="100">
            <template slot-scope="scope">
              <el-tag 
                :type="getPriorityType(scope.row.priority)" 
                size="small"
              >
                {{ getPriorityLabel(scope.row.priority) }}
              </el-tag>
            </template>
          </el-table-column>
          
          <el-table-column prop="dueDate" label="截止日期" width="120">
            <template slot-scope="scope">
              <span :class="{ 'overdue': isOverdue(scope.row.dueDate) && !scope.row.completed }">
                {{ formatDate(scope.row.dueDate, 'date') || '无' }}
              </span>
            </template>
          </el-table-column>
          
          <el-table-column prop="createdAt" label="创建时间" width="120">
            <template slot-scope="scope">
              {{ formatDate(scope.row.createdAt) }}
            </template>
          </el-table-column>
          
          <el-table-column label="操作" width="150" align="center">
            <template slot-scope="scope">
              <el-button
                size="mini"
                type="primary"
                icon="el-icon-edit"
                @click="editTask(scope.row)"
                circle
              ></el-button>
              <el-button
                size="mini"
                type="danger"
                icon="el-icon-delete"
                @click="confirmDeleteTask(scope.row)"
                circle
              ></el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-main>
    </el-container>
    
    <!-- 添加/编辑任务对话框 -->
    <el-dialog :title="dialogTitle" :visible.sync="dialogVisible" width="500px">
      <el-form :model="taskForm" :rules="taskRules" ref="taskForm" label-width="80px">
        <el-form-item label="任务名称" prop="title">
          <el-input v-model="taskForm.title" placeholder="请输入任务名称"></el-input>
        </el-form-item>
        
        <el-form-item label="任务描述" prop="description">
          <el-input
            type="textarea"
            v-model="taskForm.description"
            placeholder="请输入任务描述"
            :rows="4"
          ></el-input>
        </el-form-item>
        
        <el-form-item label="优先级" prop="priority">
          <el-select v-model="taskForm.priority" placeholder="请选择优先级">
            <el-option label="低" value="low"></el-option>
            <el-option label="中" value="medium"></el-option>
            <el-option label="高" value="high"></el-option>
          </el-select>
        </el-form-item>
        
        <el-form-item label="截止日期" prop="dueDate">
          <el-date-picker
            v-model="taskForm.dueDate"
            type="datetime"
            placeholder="选择截止日期时间"
            format="yyyy-MM-dd HH:mm"
            value-format="yyyy-MM-dd HH:mm:ss"
            :picker-options="{firstDayOfWeek: 1}"
          ></el-date-picker>
        </el-form-item>
        
        <el-form-item label="状态" prop="completed">
          <el-switch v-model="taskForm.completed" active-text="已完成" inactive-text="未完成"></el-switch>
        </el-form-item>
      </el-form>
      
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitTaskForm" :loading="submitting">确定</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'

export default {
  name: 'Home',
  data() {
    return {
      // 加载状态
      loading: false,
      // 提交状态
      submitting: false,
      // 任务过滤器
      taskFilter: 'all',
      // 优先级过滤器
      priorityFilter: 'all',
      // 对话框可见性
      dialogVisible: false,
      // 对话框标题
      dialogTitle: '新建任务',
      // 是否为编辑模式
      isEdit: false,
      // 任务表单
      taskForm: {
        id: null,
        title: '',
        description: '',
        priority: 'medium',
        dueDate: null,
        completed: false
      },
      // 表单验证规则
      taskRules: {
        title: [
          { required: true, message: '请输入任务名称', trigger: 'blur' },
          { min: 1, max: 50, message: '任务名称长度应为1-50个字符', trigger: 'blur' }
        ]
      }
    }
  },
  computed: {
    ...mapGetters({
      user: 'getUser',
      tasks: 'getTasks'
    }),
    
    // 获取用户名首字母（无头像时显示）
    userInitials() {
      if (!this.user || !this.user.username) return '?'
      return this.user.username.charAt(0).toUpperCase()
    },
    // 根据过滤条件筛选任务
    filteredTasks() {
      let filtered = this.tasks;
      
      // 按完成状态筛选
      if (this.taskFilter === 'active') {
        filtered = filtered.filter(task => !task.completed);
      } else if (this.taskFilter === 'completed') {
        filtered = filtered.filter(task => task.completed);
      }
      
      // 按优先级筛选
      if (this.priorityFilter !== 'all') {
        filtered = filtered.filter(task => task.priority === this.priorityFilter);
      }
      
      return filtered;
    }
  },
  created() {
    // 获取用户信息和任务列表
    this.fetchData()
  },
  methods: {
    // 获取数据
    async fetchData() {
      this.loading = true
      try {
        // 如果没有用户信息，获取用户信息
        if (!this.user) {
          await this.$store.dispatch('fetchUserInfo')
        }
        // 获取任务列表
        await this.$store.dispatch('fetchTasks')
      } catch (error) {
        this.$message.error('获取数据失败')
        console.error(error)
      } finally {
        this.loading = false
      }
    },
    
    // 格式化日期
    formatDate(dateString, type = 'datetime') {
      if (!dateString) return ''
      const date = new Date(dateString)
      
      if (type === 'date') {
        return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
      }
      
      return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
    },
    
    // 获取优先级标签
    getPriorityLabel(priority) {
      const labels = {
        low: '低',
        medium: '中',
        high: '高'
      }
      return labels[priority] || '中'
    },
    
    // 获取优先级标签类型
    getPriorityType(priority) {
      const types = {
        low: 'info',
        medium: 'warning',
        high: 'danger'
      }
      return types[priority] || 'warning'
    },
    
    // 判断是否过期
    isOverdue(dueDate) {
      if (!dueDate) return false
      return new Date(dueDate) < new Date()
    },
    
    // 显示添加任务对话框
    showAddTaskDialog() {
      this.isEdit = false
      this.dialogTitle = '新建任务'
      this.resetTaskForm()
      this.dialogVisible = true
    },
    
    // 编辑任务
    editTask(task) {
      this.isEdit = true
      this.dialogTitle = '编辑任务'
      this.taskForm = {
        id: task.id,
        title: task.title,
        description: task.description,
        priority: task.priority,
        dueDate: task.dueDate,
        completed: task.completed
      }
      this.dialogVisible = true
    },
    
    // 重置任务表单
    resetTaskForm() {
      if (this.$refs.taskForm) {
        this.$refs.taskForm.resetFields()
      }
      this.taskForm = {
        id: null,
        title: '',
        description: '',
        priority: 'medium',
        dueDate: null,
        completed: false
      }
    },
    
    // 提交任务表单
    submitTaskForm() {
      this.$refs.taskForm.validate(async valid => {
        if (!valid) return
        
        this.submitting = true
        
        try {
          if (this.isEdit) {
            // 更新任务
            await this.$store.dispatch('updateTask', {
              id: this.taskForm.id,
              taskData: {
                title: this.taskForm.title,
                description: this.taskForm.description,
                priority: this.taskForm.priority,
                dueDate: this.taskForm.dueDate,
                completed: this.taskForm.completed
              }
            })
            this.$message.success('任务更新成功')
          } else {
            // 创建任务
            await this.$store.dispatch('createTask', {
              title: this.taskForm.title,
              description: this.taskForm.description,
              priority: this.taskForm.priority,
              dueDate: this.taskForm.dueDate,
              completed: this.taskForm.completed
            })
            this.$message.success('任务创建成功')
          }
          
          // 关闭对话框
          this.dialogVisible = false
        } catch (error) {
          this.$message.error(this.isEdit ? '更新任务失败' : '创建任务失败')
          console.error(error)
        } finally {
          this.submitting = false
        }
      })
    },
    
    // 更新任务状态
    async updateTaskStatus(task) {
      try {
        await this.$store.dispatch('updateTask', {
          id: task.id,
          taskData: {
            title: task.title,
            description: task.description,
            priority: task.priority,
            dueDate: task.dueDate,
            completed: task.completed
          }
        })
      } catch (error) {
        this.$message.error('更新任务状态失败')
        console.error(error)
        // 恢复原状态
        task.completed = !task.completed
      }
    },
    
    // 确认删除任务
    confirmDeleteTask(task) {
      this.$confirm('确定要删除这个任务吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(async () => {
        try {
          await this.$store.dispatch('deleteTask', task.id)
          this.$message.success('任务删除成功')
        } catch (error) {
          this.$message.error('删除任务失败')
          console.error(error)
        }
      }).catch(() => {
        // 取消删除，不做任何操作
      })
    },
    
    // 登出
    logout() {
      this.$confirm('确定要退出登录吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        this.$store.dispatch('logout')
        this.$router.push('/login')
        this.$message.success('已退出登录')
      }).catch(() => {
        // 取消登出，不做任何操作
      })
    }
  }
}
</script>

<style scoped>
.home-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.el-header {
  background-color: #2c3e50;
  color: white;
  line-height: 60px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 15px;
}

.logo-section {
  flex: 1;
}

.logo-section h2 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: #fff;
}

.nav-links {
  display: flex;
  align-items: center;
  flex: 2;
  justify-content: center;
}

.nav-links .el-button {
  color: #e0e0e0;
  font-size: 15px;
  margin: 0 15px;
  padding: 10px 15px;
  border-radius: 4px;
  transition: all 0.3s;
}

.nav-links .el-button:hover {
  color: white;
  background-color: rgba(255, 255, 255, 0.1);
}

.user-info {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  flex: 1;
}

.username-container {
  display: flex;
  align-items: center;
  background-color: rgba(255, 255, 255, 0.1);
  padding: 5px 12px;
  border-radius: 20px;
  margin-right: 15px;
  cursor: pointer;
  transition: all 0.3s;
}

.username-container:hover {
  background-color: rgba(255, 255, 255, 0.2);
}

.username-text {
  margin-left: 8px;
  font-weight: 500;
}

.avatar-container {
  position: relative;
}

.mini-avatar {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid rgba(255, 255, 255, 0.8);
}

.avatar-placeholder {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  background-color: #e74c3c;
  color: white;
  display: flex;
  justify-content: center;
  align-items: center;
  font-size: 14px;
  font-weight: bold;
  border: 2px solid rgba(255, 255, 255, 0.8);
}

.logout-btn {
  font-weight: 500;
}

.el-main {
  padding: 20px;
  background-color: #f5f7fa;
}

.task-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.task-filter {
  margin-bottom: 20px;
  display: flex;
  flex-wrap: wrap;
  gap: 15px;
}

.filter-row {
  display: flex;
  align-items: center;
}

.filter-label {
  margin-right: 10px;
  font-weight: bold;
}

.overdue {
  color: #F56C6C;
  font-weight: bold;
}

.task-completed {
  text-decoration: line-through;
  color: #909399;
}
</style>
