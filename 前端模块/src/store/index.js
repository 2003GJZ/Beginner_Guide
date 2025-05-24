import Vue from 'vue'
import Vuex from 'vuex'
import axios from 'axios'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    // 用户信息
    user: null,
    // 任务列表
    tasks: []
  },
  mutations: {
    // 设置用户信息
    setUser(state, user) {
      state.user = user
    },
    // 设置任务列表
    setTasks(state, tasks) {
      state.tasks = tasks
    },
    // 添加新任务
    addTask(state, task) {
      state.tasks.push(task)
    },
    // 更新任务
    updateTask(state, updatedTask) {
      const index = state.tasks.findIndex(task => task.id === updatedTask.id)
      if (index !== -1) {
        state.tasks.splice(index, 1, updatedTask)
      }
    },
    // 删除任务
    deleteTask(state, taskId) {
      state.tasks = state.tasks.filter(task => task.id !== taskId)
    }
  },
  actions: {
    // 登录
    async login({ commit }, credentials) {
      try {
        console.log('发送登录请求数据:', credentials);
        
        // 确保密码字段正确传递
        const data = JSON.stringify({
          username: credentials.username,
          password: credentials.password
        });
        
        const config = {
          headers: {
            'Content-Type': 'application/json'
          }
        };
        
        const response = await axios.post('/api/login', data, config);
        console.log('登录响应:', response.data);
        // 保存token到本地存储
        localStorage.setItem('token', response.data.token);
        // 保存用户信息到状态
        commit('setUser', response.data.user);
        return response;
      } catch (error) {
        console.error('登录错误:', error.response ? error.response.data : error.message);
        throw error;
      }
    },
    // 注册
    async register(_, userData) {
      try {
        console.log('发送注册请求数据:', userData);
        
        // 确保密码字段正确传递
        const data = JSON.stringify({
          username: userData.username,
          password: userData.password
        });
        
        const config = {
          headers: {
            'Content-Type': 'application/json'
          }
        };
        
        const response = await axios.post('/api/register', data, config);
        console.log('注册响应:', response.data);
        return response;
      } catch (error) {
        console.error('注册错误:', error.response ? error.response.data : error.message);
        throw error;
      }
    },
    // 获取用户信息
    async fetchUserInfo({ commit }) {
      try {
        const response = await axios.get('/api/user/info')
        commit('setUser', response.data)
        return response
      } catch (error) {
        throw error
      }
    },
    // 获取任务列表
    async fetchTasks({ commit }) {
      try {
        const response = await axios.get('/api/tasks')
        commit('setTasks', response.data)
        return response
      } catch (error) {
        throw error
      }
    },
    // 创建新任务
    async createTask({ commit }, taskData) {
      try {
        const response = await axios.post('/api/task', taskData)
        commit('addTask', response.data)
        return response
      } catch (error) {
        throw error
      }
    },
    // 更新任务
    async updateTask({ commit }, { id, taskData }) {
      try {
        // 使用POST请求替代PUT，与后端路由保持一致
        const response = await axios.post(`/api/task/update/${id}`, taskData)
        commit('updateTask', response.data)
        return response
      } catch (error) {
        throw error
      }
    },
    // 删除任务
    async deleteTask({ commit }, id) {
      try {
        // 使用POST请求替代DELETE，与后端路由保持一致
        await axios.post(`/api/task/delete/${id}`)
        commit('deleteTask', id)
      } catch (error) {
        throw error
      }
    },
    // 登出
    logout({ commit }) {
      // 清除本地存储的token
      localStorage.removeItem('token')
      // 清除用户信息
      commit('setUser', null)
      // 清除任务列表
      commit('setTasks', [])
    }
  },
  getters: {
    // 获取用户信息
    getUser: state => state.user,
    // 获取任务列表
    getTasks: state => state.tasks,
    // 获取已完成任务
    getCompletedTasks: state => state.tasks.filter(task => task.completed),
    // 获取未完成任务
    getIncompleteTasks: state => state.tasks.filter(task => !task.completed)
  }
})
