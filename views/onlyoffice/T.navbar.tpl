{{define "navbar"}}
<!-- navbar-inverse一个带有黑色背景白色文本的导航栏 
固定在页面的顶部，向 .navbar class 添加 class .navbar-fixed-top
为了防止导航栏与页面主体中的其他内容
的顶部相交错，需要向 <body> 标签添加内边距，内边距的值至少是导航栏的高度。
-->
<style type="text/css">
a.navbar-brand {
  display: none;
}

@media (max-width: 960px) {
  a.navbar-brand {
    display: inline-block;
  }
}

#NavmodalTable .modal-header {
  cursor: move;
}
</style>
<link rel="stylesheet" type="text/css" href="/static/font-awesome-4.7.0/css/font-awesome.min.css" />
<nav class="navbar navbar-default navbar-static-top" style="margin-bottom: 5px;" role="navigation">
  <div class="navbar-header">
    <button type="button" class="navbar-toggle" data-toggle="collapse" data-target="#target-menu">
      <span class="sr-only">3xxx</span>
      <span class="icon-bar"></span>
      <span class="icon-bar"></span>
      <span class="icon-bar"></span>
    </button>
    <a id="11" class="navbar-brand">EngineerCMS</a>
  </div>
  <div class="collapse navbar-collapse" id="target-menu">
    <ul class="nav navbar-nav">
      <li {{if .IsIndex}} class="active" {{end}}>
        <a href="/index">首页</a>
      </li>
      <li {{if .IsProject}} class="active" {{end}}>
        <a href="/project/" id="project">项目</a>
      </li>
      <!-- **********定制导航条菜单开始******** -->
      
      <!-- **********定制导航条菜单结束******** -->
      <li {{if .IsOnlyOffice}} class="active" {{end}}>
        <a href="/onlyoffice">OnlyOffice</a>
      </li>

    </ul>
    <div class="pull-right">
      <ul class="nav navbar-nav navbar-right">
        <li><a class="navbar-brand" href="javascript:void(0)" onclick="chooseProjectButton()" id="chooseProject">切换项目</a></li>
        {{if eq true .IsAdmin}}
        <li class="dropdown">
          <a href="javascript:void(0)" class="dropdown-toggle" data-toggle="dropdown">{{.Username}} <b class="caret"></b></a>
          <ul class="dropdown-menu">
            <li><a href="/user" title="个人中心">个人中心</a></li>
            <li><a href="/admin" title="管理">进入后台</a></li>
            <li><a href="javascript:void(0)" id="login">重新登录</a></li>
            <li><a href="/v1/wx/ssologin" title="单点登录">SSO单点登陆</a></li>
            <li><a href="javascript:void(0)" onclick="logout()">退出</a></li>
          </ul>
        </li>
        {{else if eq true .IsLogin}}
        <li class="dropdown">
          <a href="javascript:void(0)" class="dropdown-toggle" data-toggle="dropdown">{{.Username}} <b class="caret"></b></a>
          <ul class="dropdown-menu">
            <li><a href="/user" title="个人中心">个人中心</a></li>
            <li><a href="javascript:void(0)" id="login">重新登录</a></li>
            <li><a href="javascript:void(0)" onclick="logout()">退出</a></li>
          </ul>
        </li>
        {{else}}
        <li class="dropdown">
          <a href="javascript:void(0)" class="dropdown-toggle" data-toggle="dropdown">{{.Username}} <b class="caret"></b></a>
          <ul class="dropdown-menu">
            <li><a href="/v1/wx/wxlogin" title="微信扫码登录">微信扫码登陆</a></li>
            <li><a href="javascript:void(0)" id="login">用户名密码登陆</a></li>
            <li><a href="/v1/wx/ssologin" title="单点登录">SSO单点登陆</a></li>
          </ul>
        </li>
        {{end}}
      </ul>
    </div>
  </div>
</nav>
<!-- 登录模态框 -->
<div class="form-horizontal">
  <div class="modal fade" id="modalNav">
    <div class="modal-dialog" id="modalDialog">
      <div class="modal-content">
        <div class="modal-header" style="background-color: #8bc34a">
          <button type="button" class="close" data-dismiss="modal">
            <span aria-hidden="true">&times;</span>
          </button>
          <h3 class="modal-title">登录</h3>
          <label id="status"></label>
        </div>
        <div class="modal-body">
          <div class="modal-body-content">
            <div class="form-group" style="width: 100%;">
                <label class="col-sm-3 control-label">用户名</label>
              <div class="col-sm-7">
                <input id="username" name="username" type="text" value="" class="form-control" placeholder="Enter Account" list="cars" onkeypress="getKey()"></div>
            </div>
            <div class="form-group" style="width: 100%;">
              <label class="col-sm-3 control-label">密码</label>
              <div class="col-sm-7">
                <input id="password" name="password" type="password" value="" class="form-control" placeholder="Enter Password"  autocomplete="off" onkeypress="getKey()"></div>
            </div>
            <div class="form-group" style="width: 100%;">
              <label class="col-sm-3 control-label"><input type="checkbox">自动登陆</label>
            </div>
            <div class="form-group" style="width: 100%;">
              <label class="col-sm-3 control-label"><a href="/regist">  没有账号？注册</a></label>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-default" data-dismiss="modal">关闭</button>
          <button type="button" class="btn btn-primary" id="submit" onclick="handleLogin()">登录</button>
        </div>
      </div>
    </div>
  </div>
</div>
<!-- 项目切换模态框 -->
<div class="form-horizontal">
  <div class="modal fade" id="NavmodalTable">
    <div class="modal-dialog">
      <div class="modal-content">
        <div class="modal-header">
          <button type="button" class="close" data-dismiss="modal">
            <span aria-hidden="true">&times;</span>
          </button>
          <h3 class="modal-title">切换项目</h3>
        </div>
        <div class="modal-body">
          <table id="Navtable2"></table>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
          <button type="button" class="btn btn-primary" id="saveproj" onclick="setlocalstorage()">保存</button>
        </div>
      </div>
    </div>
  </div>
</div>
<!-- <script type="text/javascript" src="/static/js/jquery-3.3.1.min.js"></script> -->
<script type="text/javascript" src="/static/js/jsencrypt.min.js"></script>
<script type="text/javascript">
  // $(document).keydown(function(event){
  //   if(event.keyCode==13){
  //     document.getElementById("submit").click();
  //   }
  // })
  function getKey() {
    if (event.keyCode == 13) {
      handleLogin()
    }
  }

  let encryptor = new JSEncrypt();

  // 获取公钥
  async function getPublicKey() {
    const response = await fetch('/v1/wx/publickeyhandler');
    return await response.text();
  }

  // 加密处理
  async function encryptData(data) {
    const publicKey = await getPublicKey();
    encryptor.setPublicKey(publicKey);
    return encryptor.encrypt(data);
  }

  // 弹出登录模态框
  $("#login").click(function() {
    $('#modalNav').modal({
      show: true,
      backdrop: 'static'
    });
  })
  // 登录处理
  async function handleLogin() {
    const username = document.getElementById('username').value;
    if (username.length == 0) {
      alert("请输入账号");
      return
    }

    const password = document.getElementById('password').value;
    if (password.length == 0) {
      alert("请输入密码");
      return
    }
    try {
      const encryptedUsername = await encryptData(username);
      const encryptedPassword = await encryptData(password);

      const response = await fetch('/v1/wx/loginpost', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          username: encryptedUsername,
          password: encryptedPassword
        })
      });

      const result = await response.json();
      // alert(result.msg);
      if (result.islogin == 1) {
          $("#status").html(result.msg);
          $('#modalNav').modal('hide');
          window.location.reload();
        } else if (result.islogin == 0) {
          $("#status").html(result.msg)
        } else if (result.islogin == 2) {
          $("#status").html(result.msg)
        }
    } catch (error) {
      console.error('登录失败:', error);
      alert('登录失败，请检查控制台');
    }
  }

  $(function() {
    var projectid = window.localStorage.getItem('projectid')
    if (projectid == null) {
      document.getElementById('chooseProject').innerText = '选择项目';
    } else {
  
    }
  })
  
  //登陆功能_作废！
  function login_back() {
    var uname = document.getElementById("uname");
    if (uname.value.length == 0) {
      alert("请输入账号");
      return
    }
    var pwd = document.getElementById("pwd");
    if (pwd.value.length == 0) {
      alert("请输入密码");
      return
    }
  
    $.ajax({
      type: 'post',
      url: '/loginpost',
      data: {
        "uname": $("#uname").val(),
        "pwd": $("#pwd").val()
      },
      success: function(result) {
        if (result.islogin == 1) {
          $("#status").html("登陆成功");
          $('#modalNav').modal('hide');
          window.location.reload();
        } else if (result.islogin == 0) {
          $("#status").html("用户名或密码错误！")
        } else if (result.islogin == 2) {
          $("#status").html("密码错误")
        }
      }
    })
  }
  //登出功能
  function logout() {
    $.ajax({
      type: 'get',
      url: '/logout',
      data: {},
      success: function(result) {
        if (!result.islogin) {
          alert("登出成功");
          window.location.reload();
        } else {
          alert("登出失败")
        }
      }
    })
  }
  
  $(function() {
    var projectid = window.localStorage.getItem('projectid')
    // 初始化【未接受】工作流表格
    $("#Navtable2").bootstrapTable({
      url: '/project/getprojects',
      method: 'get',
      search: 'true',
      showRefresh: 'true',
      showToggle: 'true',
      showColumns: 'true',
      // toolbar:'#toolbar1',
      pagination: 'true',
      sidePagination: "server",
      queryParamsType: '',
      //请求服务器数据时，你可以通过重写参数的方式添加一些额外的参数，例如 toolbar 中的参数 如果 queryParamsType = 'limit' ,返回参数必须包含
      // limit, offset, search, sort, order 否则, 需要包含:
      // pageSize, pageNumber, searchText, sortName, sortOrder.
      // 返回false将会终止请求。
      pageSize: 15,
      pageNumber: 1,
      pageList: [15, 20, 50, 100],
      singleSelect: "true",
      clickToSelect: "true",
      selectItemName: "project",
      queryParams: function queryParams(params) { //设置查询参数
        var param = {
          limit: params.pageSize, //每页多少条数据
          pageNo: params.pageNumber, // 页码
          searchText: params.searchText // $(".search .form-control").val()
        };
        //搜索框功能
        //当查询条件中包含中文时，get请求默认会使用ISO-8859-1编码请求参数，在服务端需要对其解码
        // if (null != searchText) {
        //   try {
        //     searchText = new String(searchText.getBytes("ISO-8859-1"), "UTF-8");
        //   } catch (Exception e) {
        //     e.printStackTrace();
        //   }
        // }
        return param;
      },
      columns: [{
          title: '选择',
          radio: 'true',
          width: '10',
          align: "center",
          valign: "middle",
          formatter: function(value, row, index) {
            return {checked: row.Id==projectid}
          },
        },
        {
          // field: 'Number',
          title: '序号',
          formatter: function(value, row, index) {
            return index + 1
          },
          align: "center",
          valign: "middle"
        },
        {
          field: 'Code',
          title: '编号',
          // formatter:setCode,
          align: "center",
          valign: "middle"
        },
        {
          field: 'Title',
          title: '名称',
          // formatter:setTitle,
          align: "center",
          valign: "middle"
        },
        // {
        //   field: 'Label',
        //   title: '标签',
        //   formatter:setLable,
        //   align:"center",
        //   valign:"middle"
        // },
        {
          field: 'Principal',
          title: '负责人',
          align: "center",
          valign: "middle"
        },
        // {
        //   field: 'Number',
        //   title: '成果数量',
        //   formatter:setCode,
        //   align:"center",
        //   valign:"middle"
        // },
        // {
        //   field: 'action',
        //   title: '时间轴',
        //   formatter:actionFormatter,
        //   events:actionEvents,
        //   align:"center",
        //   valign:"middle"
        // },
        {
          field: 'Created',
          title: '建立时间',
          formatter: localDateFormatter,
          align: "center",
          valign: "middle"
        }
        // {
        // field: 'dContMainEntity.createTime',
        // title: '发起时间',
        // formatter: function (value, row, index) {
        // return new Date(value).toLocaleString().substring(0,9);
        // }
        // },
        // {
        // field: 'dContMainEntity.operate',
        // title: '操作',
        // formatter: operateFormatter
        // }
      ]
    });
  });
  
  // 切换项目
  function chooseProjectButton() {
    if (!{{.IsLogin }}) {
      alert("请登录！");
      return;
    }
  
    // $('#project').$("input[name='types']").attr('checked','true');
  
    $('#NavmodalTable').modal({
      show: true,
      backdrop: 'static'
    });
  }
  
  // 将选择的项目id存入浏览器内存
  function setlocalstorage() {
    var selectRow2 = $('#Navtable2').bootstrapTable('getSelections');
    if (selectRow2.length < 1) {
      alert("请先勾选项目！");
      return;
    }
    console.log(selectRow2[0].Id)
    window.localStorage.setItem('projectid', selectRow2[0].Id);
    $('#NavmodalTable').modal('hide');
    // window.location.reload();
  }
  
  // $(function() {
  //   $("#NavmodalTable").draggable({ handle: ".modal-header" }); //为模态对话框添加拖拽
  //   $("#myModal").css("overflow", "hidden"); //禁止模态对话框的半透明背景滚动
  // })
  
  function localDateFormatter(value) {
    return moment(value, 'YYYY-MM-DD').format('YYYY-MM-DD');
  }
  </script>
{{end}}
