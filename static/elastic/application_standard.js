/* globals Vue */

(function() {
  'use strict'

  var app = new Vue({
    el: '#app',

    data: {
      state: {
        error: null,
        loading: true,
        fetching: false,
        replaceResults: true
      },

      isBottom: false,
      isAdmin: false,
      showLoginDialog: false,
      loginForm: {
        username: '',
        password: ''
      },

      query: '',
      results: [],
      total: 0,
      search_after: []
    },

    methods: {
      loadResults: function() {
        var self = this

        if (self.state.fetching) { return }
        self.state.fetching = true

        if (self.query === '') {
          window.history.pushState({}, '', '?')
        } else {
          window.history.pushState({}, '', '?q=' + encodeURIComponent(self.query))
        }

        if (self.state.replaceResults) { self.search_after = [] }

        // window.fetch(`/v1/elastic/search?q=${encodeURIComponent(self.query)}&a=${self.search_after.join(',')}`)
        window.fetch(`/v1/wx/elasticsearch?q=${encodeURIComponent(self.query)}&a=${self.search_after.join(',')}`)
          .then(function(response) {
            if (!response.ok) { return Promise.reject(response) }
            return response.json()
          })
          .then(function(response) {
            var results = []

            response.hits.forEach(function(r) {
              var author = {
                first_name: r.author.first_name
              }
              var result = {
                id: r.id,
                url: r.url,//+'#page=3&keyword='+${encodeURIComponent(self.query)}
                image_url: r.image_url,
                published: r.published,
                body: r.body,
                author: author
              }

              if (r.highlights && r.highlights.title) {
                result.title = r.highlights.title[0]
              } else {
                result.title = r.title
              }

              if (r.highlights && r.highlights.alt) {
                result.alt = r.highlights.alt[0]
              }

              if (r.highlights && r.highlights.body) {
                // 获得<em>位置——截取前100和后100个字符，不满100按实际算
                var myString=""; //= r.highlights.body[0];
                r.highlights.body.forEach((elem, index) => {
                  console.log(elem, index);
                  myString=myString+'……'+elem;
                });

                // var myString = r.highlights.body[0];
                // var w = myString.indexOf("<em>");
                // console.log(w)
                // var start = 0
                // if (Number(w) >= 100) {
                //   start = w - 100
                // } else {
                //   start = w
                // }
                // var myString2 = myString.substr(start, 318)
                // if (myString2.length>317){
                //   myString=myString.substr(start, 317)+"......"
                // }
                result.body = myString//r.highlights.body[0]
              }

              if (r.highlights && r.highlights.transcript) {
                result.transcript = r.highlights.transcript.join('&hellip;')
              }

              results.push(result)
            })

            self.total = response.total

            if (self.state.replaceResults) {
              self.results = results
            } else {
              self.results = self.results.concat(results)
            }

            if (response.hits.length > 0) {
              self.search_after = response.hits[response.hits.length - 1].sort
            }
          })
          .then(function() {
            self.state.loading = false
            self.state.fetching = false
          })
          .catch(function(response) {
            self.state.loading = false
            self.state.fetching = false
            self.state.error = response
          })
      },

      toggle: function(event) {
        event.currentTarget.closest('div.result').classList.toggle('expanded')
      },
      
      toggleAdminMode: function() {
        var self = this
        
        // 如果已经是管理员模式，直接退出
        if (this.isAdmin) {
          this.logoutAdmin()
          return
        }
        
        // 显示登录对话框
        this.showLoginDialog = true
      },
      
      adminLogin: function() {
        var self = this
        
        // 发送登录请求到服务端API
        window.fetch('/v1/wx/admin/login', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            username: this.loginForm.username,
            password: this.loginForm.password
          })
        })
        .then(function(response) {
          if (!response.ok) {
            return Promise.reject(response)
          }
          return response.json()
        })
        .then(function(data) {
          // 登录成功
          if (data.success) {
            self.isAdmin = true
            self.showLoginDialog = false
            self.loginForm.username = ''
            self.loginForm.password = ''
            console.log('管理员登录成功')
            alert('管理员登录成功！')
          } else {
            alert('登录失败：' + (data.message || '用户名或密码错误'))
          }
        })
        .catch(function(error) {
          console.error('登录失败:', error)
          
          // 如果后端API不支持管理员登录，提供备选方案
          if (error.status === 404) {
            alert('后端管理员API暂不可用，使用本地验证模式')
            // 本地验证（仅用于演示，生产环境应使用服务端验证）
            if (self.loginForm.username === 'admin' && self.loginForm.password === 'admin123') {
              self.isAdmin = true
              self.showLoginDialog = false
              self.loginForm.username = ''
              self.loginForm.password = ''
              alert('管理员登录成功（本地验证模式）！')
            } else {
              alert('用户名或密码错误（默认：admin/admin123）')
            }
          } else {
            alert('登录失败，请检查网络连接或联系系统管理员')
          }
        })
      },
      
      cancelLogin: function() {
        this.showLoginDialog = false
        this.loginForm.username = ''
        this.loginForm.password = ''
      },
      
      logoutAdmin: function() {
        var self = this
        
        // 发送退出登录请求
        window.fetch('/v1/wx/logout', {
          method: 'get'
        })
        .then(function(response) {
          // 无论成功与否都退出管理员模式
          self.isAdmin = false
          console.log('退出管理员模式')
          alert('已退出管理员模式')
        })
        .catch(function(error) {
          // 即使API调用失败也退出管理员模式
          self.isAdmin = false
          console.log('退出管理员模式')
          alert('已退出管理员模式')
        })
      },
      
      checkAdminSession: function() {
        var self = this
        
        // 检查本地存储中是否有管理员会话
        // var adminSession = localStorage.getItem('elasticsearch_admin_session')
        // if (adminSession) {
          // 验证会话是否有效
          // window.fetch('/v1/wx/admin/check-session', {
          window.fetch('/v1/wx/islogin', {
            // headers: {
            //   'Authorization': 'Bearer ' + adminSession
            // }
          })
          .then(function(response) {
            if (response.ok) {
              return response.json()
            }
            return Promise.reject(response)
          })
          .then(function(data) {
            // if (data.valid) {
            if (data.isadmin) {
              self.isAdmin = true
              console.log('管理员会话有效')
            }
          })
          .catch(function(error) {
            // 会话无效，清除本地存储
            localStorage.removeItem('elasticsearch_admin_session')
            console.log('管理员会话已过期')
          })
        // }
      },
      
      deleteResult: function(resultId) {
        var self = this
        
        // 确认删除
        if (!confirm('确定要删除这条记录吗？此操作不可撤销。')) {
          return
        }
        
        // 发送删除请求到后端API
        window.fetch(`/v1/wx/deleteelasticsearch/${resultId}`, {
          method: 'DELETE',
          headers: {
            'Content-Type': 'application/json'
          }
        })
        .then(function(response) {
          if (!response.ok) {
            return Promise.reject(response)
          }
          return response.json()
        })
        .then(function(response) {
          // 从前端移除已删除的记录
          self.results = self.results.filter(function(result) {
            return result.id !== resultId
          })
          self.total = self.total - 1
          
          // 显示成功消息
          alert('记录删除成功！')
        })
        .catch(function(error) {
          console.error('删除失败:', error)
          alert('删除失败，请稍后重试。')
          
          // 如果后端API不支持删除，提供备选方案
          if (error.status === 404 || error.status === 405) {
            alert('后端API暂不支持删除功能，已从前端临时移除该记录。刷新页面后会恢复。')
            // 临时从前端移除记录
            self.results = self.results.filter(function(result) {
              return result.id !== resultId
            })
            self.total = self.total - 1
          }
        })
      },

      deleteElasticAll: function(resultId) {
        var self = this
        
        // 确认删除
        if (!confirm('确定要删除所有记录吗？此操作不可撤销。')) {
          return
        }
        
        // 发送删除请求到后端API
        window.fetch(`/v1/wx/deleteelasticall`, {
          method: 'DELETE',
          headers: {
            'Content-Type': 'application/json'
          }
        })
        .then(function(response) {
          if (!response.ok) {
            return Promise.reject(response)
          }
          return response.json()
        })
        .then(function(response) {
          // 从前端移除已删除的记录
          self.results = self.results.filter(function(result) {
            return result.id !== resultId
          })
          self.total = self.total - 1
          
          // 显示成功消息
          alert('记录删除成功！')
        })
        .catch(function(error) {
          console.error('删除失败:', error)
          alert('删除失败，请稍后重试。')
          
          // 如果后端API不支持删除，提供备选方案
          if (error.status === 404 || error.status === 405) {
            alert('后端API暂不支持删除功能，已从前端临时移除该记录。刷新页面后会恢复。')
            // 临时从前端移除记录
            self.results = self.results.filter(function(result) {
              return result.id !== resultId
            })
            self.total = self.total - 1
          }
        })
      }
    },

    watch: {
      query: function() {
        this.state.replaceResults = true
        this.loadResults()
      },

      isBottom: function() {
        if (this.total > this.results.length) {
          this.state.replaceResults = false
          this.loadResults()
        }
      }
    },

    created: function() {
      var self = this

      var q = document.location.search.split('q=')[1]
      if (q) { self.query = decodeURIComponent(q) }
      self.loadResults()

      // 检查管理员会话状态
      this.checkAdminSession()

      window.onscroll = function() {
        self.isBottom = (document.documentElement.scrollTop || document.body.scrollTop) + window.innerHeight === document.documentElement.scrollHeight
      }
    }
  })

  return app
})()

// vue.js提示Vue is not a constructor或Vue.createApp is not a function解决方法 -->
/*Vue 3*/
// Vue.createApp({
//   data() {
//     return {
//       items: [{ message: 'Foo' }, { message: 'Bar' }]
//     }
//   }
// }).mount('#array-rendering')

/*Vue 2*/
// var example1 = new Vue({
//   el: '#example-1',
//   data: {
//     items: [
//       { message: 'Foo' },
//       { message: 'Bar' }
//     ]
//   }
// })