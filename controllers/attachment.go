// 成果里的附件操作
package controllers

import (
	"encoding/json"
	"github.com/3xxx/engineercms/controllers/utils"
	"github.com/3xxx/engineercms/models"
	"github.com/beego/beego/v2/client/httplib"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/beego/beego/v2/server/web/pagination"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type AttachController struct {
	// beego.Controller
	web.Controller
}

// type AttachmentLink struct {
// 	Id    int64
// 	Title string
// 	Link  string
// }

// 取得某个成果id下的附件(除去pdf)给table
func (c *AttachController) GetAttachments() {
	id := c.Ctx.Input.Param(":id")
	// beego.Info(id)
	c.Data["Id"] = id
	var idNum int64
	var err error
	var Url string
	if id != "" {
		//id转成64为
		idNum, err = strconv.ParseInt(id, 10, 64)
		if err != nil {
			logs.Error(err)
		}
		//由成果id（后台传过来的行id）取得侧栏目录id
		prod, err := models.GetProd(idNum)
		if err != nil {
			logs.Error(err)
		}
		//由proj id取得url
		Url, _, err = GetUrlPath(prod.ProjectId)
		if err != nil {
			logs.Error(err)
		}
		// beego.Info(Url)
	} //else {

	//}
	//根据成果id取得所有附件
	Attachments, err := models.GetAttachments(idNum)
	if err != nil {
		logs.Error(err)
	}
	//对成果进行循环
	//赋予url
	//如果是一个成果，直接给url;如果大于1个，则是数组:这个在前端实现
	// http.ServeFile(ctx.ResponseWriter, ctx.Request, filePath)
	link := make([]AttachmentLink, 0)
	for _, v := range Attachments {
		// fileext := path.Ext(v.FileName)
		if path.Ext(v.FileName) != ".pdf" && path.Ext(v.FileName) != ".PDF" {
			linkarr := make([]AttachmentLink, 1)
			linkarr[0].Id = v.Id
			linkarr[0].Title = v.FileName
			// linkarr[0].Suffix=path.Ext(v.FileName)
			linkarr[0].FileSize = v.FileSize
			linkarr[0].Downloads = v.Downloads
			linkarr[0].Created = v.Created
			linkarr[0].Updated = v.Updated
			linkarr[0].Link = Url + "/" + v.FileName
			link = append(link, linkarr...)
		}
	}
	c.Data["json"] = link
	c.ServeJSON()
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	// c.Data["json"] = root
	// c.ServeJSON()
}

// 取得同步成果下的附件列表
func (c *AttachController) GetsynchAttachments() {
	id := c.GetString("id")
	// this.GetString("jsoninfo")
	site := c.GetString("site")
	// beego.Info(site)
	link := make([]AttachmentLink, 0)
	var attachmentlink []AttachmentLink
	//	通过如下接口可以设置请求的超时时间和数据读取时间
	jsonstring, err := httplib.Get(site+"project/product/providesynchattachment?id="+id).SetTimeout(100*time.Second, 30*time.Second).String() //.ToJSON(&productlink)
	if err != nil {
		logs.Error(err)
	}
	//json字符串解析到结构体，以便进行追加
	err = json.Unmarshal([]byte(jsonstring), &attachmentlink)
	if err != nil {
		logs.Error(err)
	}
	// beego.Info(productlink)
	link = append(link, attachmentlink...)

	c.Data["json"] = link //products
	c.ServeJSON()
}

// 提供给同步用的某个成果id下的附件(除去pdf)给table
func (c *AttachController) ProvideAttachments() {
	id := c.GetString("id")
	site := c.Ctx.Input.Site() + ":" + strconv.Itoa(c.Ctx.Input.Port())
	c.Data["Id"] = id
	var idNum int64
	var err error
	var Url string
	if id != "" {
		//id转成64为
		idNum, err = strconv.ParseInt(id, 10, 64)
		if err != nil {
			logs.Error(err)
		}
		//由成果id（后台传过来的行id）取得侧栏目录id
		prod, err := models.GetProd(idNum)
		if err != nil {
			logs.Error(err)
		}
		//由proj id取得url
		Url, _, err = GetUrlPath(prod.ProjectId)
		if err != nil {
			logs.Error(err)
		}
		// beego.Info(Url)
	} else {

	}
	//根据成果id取得所有附件
	Attachments, err := models.GetAttachments(idNum)
	if err != nil {
		logs.Error(err)
	}
	//对成果进行循环
	//赋予url
	//如果是一个成果，直接给url;如果大于1个，则是数组:这个在前端实现
	// http.ServeFile(ctx.ResponseWriter, ctx.Request, filePath)
	link := make([]AttachmentLink, 0)
	for _, v := range Attachments {
		// fileext := path.Ext(v.FileName)
		if path.Ext(v.FileName) != ".pdf" && path.Ext(v.FileName) != ".PDF" {
			linkarr := make([]AttachmentLink, 1)
			linkarr[0].Id = v.Id
			linkarr[0].Title = v.FileName
			linkarr[0].FileSize = v.FileSize
			linkarr[0].Downloads = v.Downloads
			linkarr[0].Created = v.Created
			linkarr[0].Updated = v.Updated
			linkarr[0].Link = site + Url + "/" + v.FileName
			link = append(link, linkarr...)
		}
	}
	c.Data["json"] = link
	c.ServeJSON()
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	// c.Data["json"] = root
	// c.ServeJSON()
}

// 取得某个成果id下的所有附件(包含pdf和文章)给table
// 用于编辑，这个要改，不要显示文章？
// 自从文章采用link后，是否可以一同删除？
func (c *AttachController) GetAllAttachments() {
	id := c.Ctx.Input.Param(":id")
	c.Data["Id"] = id
	var idNum int64
	var err error
	var Url string
	//id转成64为
	idNum, err = strconv.ParseInt(id, 10, 64)
	if err != nil {
		logs.Error(err)
	}
	//由成果id（后台传过来的行id）取得侧栏目录id
	prod, err := models.GetProd(idNum)
	if err != nil {
		logs.Error(err)
	}
	//由proj id取得url
	Url, _, err = GetUrlPath(prod.ProjectId)
	if err != nil {
		logs.Error(err)
	}
	// beego.Info(Url)

	//根据成果id取得所有附件
	Attachments, err := models.GetAttachments(idNum)
	if err != nil {
		logs.Error(err)
	}
	//对成果进行循环
	//赋予url
	//如果是一个成果，直接给url;如果大于1个，则是数组:这个在前端实现
	// http.ServeFile(ctx.ResponseWriter, ctx.Request, filePath)
	link := make([]AttachmentLink, 0)
	for _, v := range Attachments {
		// fileext := path.Ext(v.FileName)
		linkarr := make([]AttachmentLink, 1)
		linkarr[0].Id = v.Id
		linkarr[0].Title = v.FileName
		linkarr[0].FileSize = v.FileSize
		linkarr[0].Downloads = v.Downloads
		linkarr[0].Created = v.Created
		linkarr[0].Updated = v.Updated
		linkarr[0].Link = Url + "/" + v.FileName
		link = append(link, linkarr...)
	}
	//根据成果id取得所有文章——经过论证，不能列出文章，否则会导致使用文字的id删除了对应id的附件。
	// Articles, err := models.GetArticles(idNum)
	// if err != nil {
	// 	logs.Error(err)
	// }
	// for _, x := range Articles {
	// 	linkarr := make([]AttachmentLink, 1)
	// 	linkarr[0].Id = x.Id
	// 	linkarr[0].Title = x.Subtext
	// 	linkarr[0].Created = x.Created
	// 	linkarr[0].Updated = x.Updated
	// 	linkarr[0].Link = "/project/product/article/" + strconv.FormatInt(x.Id, 10)
	// 	link = append(link, linkarr...)
	// }
	c.Data["json"] = link
	c.ServeJSON()
	// c.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
	// c.Data["json"] = root
	// c.ServeJSON()
}

// @Title get wf document details
// @Description get documentdetail
// @Param dtid query string  true "The id of doctype"
// @Param docid query string  true "The id of doc"
// @Success 200 {object} models.GetProductsPage
// @Failure 400 Invalid page supplied
// @Failure 404 data not found
// @router /project/product/pdf/:id [get]
// 2.点击一个具体文档——显示详情——显示actions
// 取得某个成果id下的附件中的pdf给table
func (c *AttachController) GetPdfs() {
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", c.Ctx.Request.Header.Get("Origin"))

	id := c.Ctx.Input.Param(":id")
	// beego.Info(id)
	c.Data["Id"] = id
	var idNum int64
	var err error
	var Url string
	if id != "" {
		//id转成64为
		idNum, err = strconv.ParseInt(id, 10, 64)
		if err != nil {
			logs.Error(err)
		}
		//由成果id（后台传过来的行id）取得侧栏目录id
		prod, err := models.GetProd(idNum)
		if err != nil {
			logs.Error(err)
		}
		//由proj id取得url
		Url, _, err = GetUrlPath(prod.ProjectId)
		if err != nil {
			logs.Error(err)
		}
		// beego.Info(Url)
	} else {

	}
	//根据成果id取得所有附件
	Attachments, err := models.GetAttachments(idNum)
	if err != nil {
		logs.Error(err)
	}
	//对成果进行循环
	//赋予url
	//如果是一个成果，直接给url;如果大于1个，则是数组:这个在前端实现
	// http.ServeFile(ctx.ResponseWriter, ctx.Request, filePath)
	link := make([]AttachmentLink, 0)
	for _, v := range Attachments {
		if path.Ext(v.FileName) == ".pdf" || path.Ext(v.FileName) == ".PDF" {
			linkarr := make([]AttachmentLink, 1)
			linkarr[0].Id = v.Id
			linkarr[0].Title = v.FileName
			linkarr[0].FileSize = v.FileSize
			linkarr[0].Downloads = v.Downloads
			linkarr[0].Created = v.Created
			linkarr[0].Updated = v.Updated
			linkarr[0].Link = Url + "/" + v.FileName
			link = append(link, linkarr...)
		}
	}
	c.Data["json"] = link
	c.ServeJSON()
}

// 取得同步成果下的pdf列表
func (c *AttachController) GetsynchPdfs() {
	id := c.GetString("id")
	site := c.GetString("site")
	// beego.Info(site)
	link := make([]PdfLink, 0)
	var pdflink []PdfLink
	//	通过如下接口可以设置请求的超时时间和数据读取时间
	jsonstring, err := httplib.Get(site+"project/product/providesynchpdf?id="+id).SetTimeout(100*time.Second, 30*time.Second).String() //.ToJSON(&productlink)
	if err != nil {
		logs.Error(err)
	}
	//json字符串解析到结构体，以便进行追加
	err = json.Unmarshal([]byte(jsonstring), &pdflink)
	if err != nil {
		logs.Error(err)
	}
	// beego.Info(productlink)
	link = append(link, pdflink...)

	c.Data["json"] = link //products
	c.ServeJSON()
}

// 提供给同步用的pdf列表数据
func (c *AttachController) ProvidePdfs() {
	id := c.GetString("id")
	site := c.Ctx.Input.Site() + ":" + strconv.Itoa(c.Ctx.Input.Port())

	c.Data["Id"] = id
	var idNum int64
	var err error
	var Url string
	if id != "" {
		//id转成64为
		idNum, err = strconv.ParseInt(id, 10, 64)
		if err != nil {
			logs.Error(err)
		}
		//由成果id（后台传过来的行id）取得侧栏目录id
		prod, err := models.GetProd(idNum)
		if err != nil {
			logs.Error(err)
		}
		//由proj id取得url
		Url, _, err = GetUrlPath(prod.ProjectId)
		if err != nil {
			logs.Error(err)
		}
		// beego.Info(Url)
	} else {

	}
	//根据成果id取得所有附件
	Attachments, err := models.GetAttachments(idNum)
	if err != nil {
		logs.Error(err)
	}
	//对成果进行循环
	//赋予url
	//如果是一个成果，直接给url;如果大于1个，则是数组:这个在前端实现
	// http.ServeFile(ctx.ResponseWriter, ctx.Request, filePath)
	link := make([]AttachmentLink, 0)
	for _, v := range Attachments {
		if path.Ext(v.FileName) == ".pdf" || path.Ext(v.FileName) == ".PDF" {
			linkarr := make([]AttachmentLink, 1)
			linkarr[0].Id = v.Id
			linkarr[0].Title = v.FileName
			linkarr[0].FileSize = v.FileSize
			linkarr[0].Downloads = v.Downloads
			linkarr[0].Created = v.Created
			linkarr[0].Updated = v.Updated
			linkarr[0].Link = site + Url + "/" + v.FileName
			link = append(link, linkarr...)
		}
	}
	c.Data["json"] = link
	c.ServeJSON()
}

// 向某个侧栏id下添加成果——用于第一种批量添加一对一模式
func (c *AttachController) AddAttachment() {
	_, _, uid, isadmin, isLogin := checkprodRole(c.Ctx)
	if !isLogin {
		// route := c.Ctx.Request.URL.String()
		// c.Data["Url"] = route
		// c.Redirect("/roleerr?url="+route, 302)
		c.Data["json"] = "未登陆"
		c.ServeJSON()
		return
	}
	useridstring := strconv.FormatInt(uid, 10)
	pid := c.GetString("pid")
	//id转成64位
	idNum, err := strconv.ParseInt(pid, 10, 64)
	if err != nil {
		logs.Error(err)
	}
	//取项目本身
	category, err := models.GetProj(idNum)
	if err != nil {
		logs.Error(err)
	}

	var topprojectid int64
	if category.ParentId != 0 { //如果不是根目录
		parentidpath := strings.Replace(strings.Replace(category.ParentIdPath, "#$", "-", -1), "$", "", -1)
		parentidpath1 := strings.Replace(parentidpath, "#", "", -1)
		patharray := strings.Split(parentidpath1, "-")
		topprojectid, err = strconv.ParseInt(patharray[0], 10, 64)
		if err != nil {
			logs.Error(err)
		}
	} else {
		topprojectid = category.Id
	}

	projectuser, err := models.GetProjectUser(topprojectid)
	if err != nil {
		logs.Error(err)
	}
	var PostPermission bool
	if projectuser.Id == uid || isadmin {
		PostPermission = true
	} else {
		//2.取得侧栏目录路径——路由id
		//2.1 根据id取得路由
		var projurls string
		proj, err := models.GetProj(idNum)
		if err != nil {
			logs.Error(err)
		}
		if proj.ParentId == 0 { //如果是项目根目录
			projurls = "/" + strconv.FormatInt(proj.Id, 10)
		} else {
			// projurls = "/" + strings.Replace(proj.ParentIdPath, "-", "/", -1) + "/" + strconv.FormatInt(proj.Id, 10)
			projurls = "/" + strings.Replace(strings.Replace(proj.ParentIdPath, "#", "/", -1), "$", "", -1) + strconv.FormatInt(proj.Id, 10)
		}

		if res, _ := e.Enforce(useridstring, projurls+"/", "POST", ".1"); res {
			PostPermission = true
		} else {
			PostPermission = false
		}
	}
	if PostPermission {
		// meritbasic, err := models.GetMeritBasic()
		// if err != nil {
		// 	logs.Error(err)
		// }
		// var news string
		// var cid int64
		// var parentidpath string
		var file_path, DiskDirectory, Url string
		// var catalog models.PostMerit

		prodlabel := c.GetString("prodlabel")
		prodprincipal := c.GetString("prodprincipal")
		size := c.GetString("size")
		filesize, err := strconv.ParseInt(size, 10, 64)
		if err != nil {
			logs.Error(err)
		}
		filesize = filesize / 1000.0
		//根据pid查出项目id
		proj, err := models.GetProj(idNum)
		if err != nil {
			logs.Error(err)
		}
		//根据proj的parentIdpath——这个已经有了专门函数，下列可以简化！
		//由proj id取得url
		Url, DiskDirectory, err = GetUrlPath(proj.Id)
		if err != nil {
			logs.Error(err)
		}
		// 20200918注释掉以下部分
		// if proj.ParentIdPath != "" {
		// 	parentidpath = strings.Replace(strings.Replace(proj.ParentIdPath, "#$", "-", -1), "$", "", -1)
		// 	parentidpath1 = strings.Replace(parentidpath, "#", "", -1)
		// 	patharray := strings.Split(parentidpath1, "-")
		// 	topprojectid, err = strconv.ParseInt(patharray[0], 10, 64)
		// 	if err != nil {
		// 		logs.Error(err)
		// 	}
		// 	meritNum, err := strconv.ParseInt(patharray[0], 10, 64)
		// 	if err != nil {
		// 		logs.Error(err)
		// 	}
		// 	meritproj, err := models.GetProj(meritNum)
		// 	if err != nil {
		// 		logs.Error(err)
		// 	}
		// 	catalog.ProjectNumber = meritproj.Code
		// 	catalog.ProjectName = meritproj.Title
		// 	if len(patharray) > 1 {
		// 		meritNum1, err := strconv.ParseInt(patharray[1], 10, 64)
		// 		if err != nil {
		// 			logs.Error(err)
		// 		}
		// 		meritproj1, err := models.GetProj(meritNum1)
		// 		if err != nil {
		// 			logs.Error(err)
		// 		}
		// 		catalog.DesignStage = meritproj1.Title
		// 	}
		// 	if len(patharray) > 2 {
		// 		meritNum2, err := strconv.ParseInt(patharray[2], 10, 64)
		// 		if err != nil {
		// 			logs.Error(err)
		// 		}
		// 		meritproj2, err := models.GetProj(meritNum2)
		// 		if err != nil {
		// 			logs.Error(err)
		// 		}
		// 		catalog.Section = meritproj2.Title
		// 	}
		// 	for _, v := range patharray {
		// 		idNum, err := strconv.ParseInt(v, 10, 64)
		// 		if err != nil {
		// 			logs.Error(err)
		// 		}
		// 		proj1, err := models.GetProj(idNum)
		// 		if err != nil {
		// 			logs.Error(err)
		// 		}
		// 		if proj1.ParentId == 0 {
		// 			DiskDirectory = "./attachment/" + proj1.Code + proj1.Title
		// 			Url = "/attachment/" + proj1.Code + proj1.Title
		// 		} else {
		// 			filepath = proj1.Title
		// 			DiskDirectory = DiskDirectory + "/" + filepath
		// 			Url = Url + "/" + filepath
		// 		}
		// 	}
		// 	DiskDirectory = DiskDirectory + "/" + proj.Title //加上自身
		// 	Url = Url + "/" + proj.Title
		// } else {
		// 	DiskDirectory = "./attachment/" + proj.Code + proj.Title //加上自身
		// 	Url = "/attachment/" + proj.Title
		// 	catalog.ProjectNumber = proj.Code
		// 	catalog.ProjectName = proj.Title
		// 	topprojectid = proj.Id
		// }

		//获取上传的文件
		_, h, err := c.GetFile("file")
		if err != nil {
			logs.Error(err)
		}
		if h != nil {
			//保存附件
			// attachment = h.Filename
			// filepath = DiskDirectory + "/" + h.Filename
			// f.Close() // 关闭上传的文件，不然的话会出现临时文件不能清除的情况
			//将附件的编号和名称写入数据库
			_, filename1, filename2, _, _, _, _ := Record(h.Filename)
			// filename1, filename2 := SubStrings(attachment)
			//当2个文件都取不到filename1的时候，数据库里的tnumber的唯一性检查出错。
			if filename1 == "" {
				filename1 = filename2 //如果编号为空，则用文件名代替，否则多个编号为空导致存入数据库唯一性检查错误
			}
			//20190728改名，替换文件名中的#和斜杠
			filename2 = strings.Replace(filename2, "#", "号", -1)
			filename2 = strings.Replace(filename2, "/", "-", -1)
			FileSuffix := path.Ext(h.Filename)
			attachmentname := filename1 + filename2 + FileSuffix
			// code := filename1
			// title := filename2
			//存入成果数据库
			//如果编号重复，则不写入，只返回Id值。
			//根据id添加成果code, title, label, principal, content string, projectid int64
			prodId, err := models.AddProduct(filename1, filename2, prodlabel, prodprincipal, uid, idNum, topprojectid)
			if err != nil {
				logs.Error(err)
			}

			//20200918注释掉以下部分。成果写入postmerit表，准备提交merit*********
			// catalog.Tnumber = filename1
			// catalog.Name = filename2
			// catalog.Count = 1
			// catalog.Drawn = meritbasic.Nickname
			// catalog.Designd = meritbasic.Nickname
			// catalog.Author = meritbasic.Username
			// catalog.Drawnratio = 0.4
			// catalog.Designdratio = 0.4
			// const lll = "2006-01-02"
			// convdate := time.Now().Format(lll)
			// t1, err := time.Parse(lll, convdate) //这里t1要是用t1:=就不是前面那个t1了
			// if err != nil {
			// 	logs.Error(err)
			// }
			// catalog.Datestring = convdate
			// catalog.Date = t1
			// catalog.Created = time.Now() //.Add(+time.Duration(hours) * time.Hour)
			// catalog.Updated = time.Now() //.Add(+time.Duration(hours) * time.Hour)
			// catalog.Complex = 1
			// catalog.State = 0
			// cid, err, news = models.AddPostMerit(catalog)
			// if err != nil {
			// 	logs.Error(err)
			// } else {
			// 	link1 := Url + "/" + attachmentname             // + FileSuffix //附件链接地址

			file_name := filepath.Clean(attachmentname)
			clean_file_name := strings.TrimPrefix(filepath.Join(string(filepath.Separator), file_name), string(filepath.Separator))
			file_path = DiskDirectory + "/" + clean_file_name // + FileSuffix
			// 	_, err = models.AddCatalogLink(cid, link1)
			// 	if err != nil {
			// 		logs.Error(err)
			// 	}
			// 	data := news
			// 	c.Ctx.WriteString(data)
			// }
			//生成提交merit的清单结束*******************

			//把成果id作为附件的parentid，把附件的名称等信息存入附件数据库
			//如果附件名称相同，则覆盖上传，但数据库不追加
			_, err = models.AddAttachment(clean_file_name, filesize, 0, prodId)
			if err != nil {
				logs.Error(err)
				c.Data["json"] = map[string]interface{}{"state": "写入附件数据库错误!", "info": "写入附件数据库错误!", "data": "写入附件数据库错误!", "err": "ERROR"}
				c.ServeJSON()
			} else {
				//存入文件夹
				err = c.SaveToFile("file", file_path) //.Join("attachment", attachment)) //存文件    WaterMark(filepath)    //给文件加水印
				if err != nil {
					logs.Error(err)
					c.Data["json"] = map[string]interface{}{"state": "存储文件错误!", "info": "存储文件错误!", "data": "存储文件错误!", "err": "ERROR"}
					c.ServeJSON()
				}
				c.Data["json"] = map[string]interface{}{"state": "SUCCESS", "title": attachmentname, "original": attachmentname, "url": Url + "/" + attachmentname}
				c.ServeJSON()
			}
		}
	} else {
		c.Data["json"] = map[string]interface{}{"state": "ERROR"}
		c.ServeJSON()
	}
	// } else {
	// route := c.Ctx.Request.URL.String()
	// c.Data["Url"] = route
	// c.Redirect("/roleerr?url="+route, 302)
	// c.Redirect("/roleerr", 302)
	// return
	// }
	// c.TplName = "topic_one_add.tpl" //不加这句上传出错，虽然可以成功上传
	// c.Redirect("/topic", 302)
	// success : 0 | 1,           // 0 表示上传失败，1 表示上传成功
	//    message : "提示的信息，上传成功或上传失败及错误信息等。",
	//    url     : "图片地址"        // 上传成功时才返回
}

// @Title post wx attachment by projectid
// @Description post attachment by projectid
// @Param pid query string true "The projectid of attachment"
// @Success 200 {object} models.AddArticle
// @Failure 400 Invalid page supplied
// @Failure 404 project not found
// @router /addwxattachment [post]
// wx向某个侧栏id下添加成果——用于第一种批量添加一对一模式
func (c *AttachController) AddWxAttachment() {
	var user models.User
	var err error
	openID := c.GetSession("openID")
	if openID != nil {
		user, err = models.GetUserByOpenID(openID.(string))
		if err != nil {
			logs.Error(err)
		}
	} else {
		c.Data["json"] = map[string]interface{}{"info": "未登陆用户"}
		c.ServeJSON()
		return
	}
	//  else {
	// 	user, err = models.GetUserByUsername("admin")
	// 	if err != nil {
	// 		logs.Error(err)
	// 	}
	// }
	var parentidpath, parentidpath1 string
	var DiskDirectory, Url string

	//解析表单
	filename := c.GetString("filename")
	pid := c.GetString("pid")
	//pid转成64为
	pidNum, err := strconv.ParseInt(pid, 10, 64)
	if err != nil {
		logs.Error(err)
	}
	//根据pid查出项目id
	proj, err := models.GetProj(pidNum)
	if err != nil {
		logs.Error(err)
	}
	var topprojectid int64
	if proj.ParentIdPath != "" {
		parentidpath = strings.Replace(strings.Replace(proj.ParentIdPath, "#$", "-", -1), "$", "", -1)
		parentidpath1 = strings.Replace(parentidpath, "#", "", -1)
		patharray := strings.Split(parentidpath1, "-")
		topprojectid, err = strconv.ParseInt(patharray[0], 10, 64)
		if err != nil {
			logs.Error(err)
		}
	} else {
		topprojectid = proj.Id
	}

	//根据proj的parentIdpath——这个已经有了专门函数，下列可以简化！
	//由proj id取得url
	Url, DiskDirectory, err = GetUrlPath(pidNum)
	if err != nil {
		logs.Error(err)
	}

	//获取上传的文件
	_, h, err := c.GetFile("file")
	if err != nil {
		logs.Error(err)
	}
	if h != nil {
		//保存附件
		// f.Close() // 关闭上传的文件，不然的话会出现临时文件不能清除的情况
		//将附件的编号和名称写入数据库
		_, filename1, filename2, _, _, _, _ := Record(filename)
		// filename1, filename2 := SubStrings(attachment)
		//当2个文件都取不到filename1的时候，数据库里的tnumber的唯一性检查出错。
		if filename1 == "" {
			filename1 = filename2 //如果编号为空，则用文件名代替，否则多个编号为空导致存入数据库唯一性检查错误
		}
		//20190728改名，替换文件名中的#和斜杠
		//文件命中怎么可能出现斜杠？？？
		filename2 = strings.Replace(filename2, "#", "号", -1)
		filename2 = strings.Replace(filename2, "/", "-", -1)
		FileSuffix := path.Ext(h.Filename)
		attachmentname := filename1 + filename2 + FileSuffix
		//存入成果数据库
		//如果编号重复，则不写入，只返回Id值。
		//根据id添加成果code, title, label, principal, content string, projectid int64
		prodId, err := models.AddProduct(filename1, filename2, "", "", user.Id, pidNum, topprojectid)
		if err != nil {
			logs.Error(err)
		}
		//除了product，其余都是原始文件名，挺好的吧。
		file_name := filepath.Clean(filename)
		clean_file_name := strings.TrimPrefix(filepath.Join(string(filepath.Separator), file_name), string(filepath.Separator))
		// file_path = DiskDirectory + "/" + clean_file_path

		file_path := DiskDirectory + "/" + clean_file_name // filename

		//如果附件名称相同，则覆盖上传，但数据库不追加
		_, err = models.AddAttachment(clean_file_name, 0, 0, prodId)
		if err != nil {
			logs.Error(err)
			c.Data["json"] = map[string]interface{}{"info": "ERR"}
			c.ServeJSON()
		} else {
			//存入文件夹
			err = c.SaveToFile("file", file_path) //给文件加水印WaterMark(filepath)
			if err != nil {
				logs.Error(err)
			}
			c.Data["json"] = map[string]interface{}{"info": "SUCCESS", "title": attachmentname, "original": attachmentname, "url": Url + "/" + attachmentname}
			c.ServeJSON()
		}
	}
}

// 向服务器保存dwg文件
func (c *AttachController) SaveDwgfile() {
	id := c.GetString("id")
	//pid转成64为
	idNum, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		logs.Error(err)
		c.Data["json"] = map[string]interface{}{"info": "ERROR", "state": "ERROR", "data": "字符转int64错误！"}
		c.ServeJSON()
		return
	}
	//根据附件id取得附件的prodid，路径
	attachment, err := models.GetAttachbyId(idNum)
	if err != nil {
		logs.Error(err)
		c.Data["json"] = map[string]interface{}{"info": "ERROR", "state": "ERROR", "data": "查询附件错误！"}
		c.ServeJSON()
		return
	}

	product, err := models.GetProd(attachment.ProductId)
	if err != nil {
		logs.Error(err)
		c.Data["json"] = map[string]interface{}{"info": "ERROR", "state": "ERROR", "data": "查询成果错误！"}
		c.ServeJSON()
		return
	}
	//由proj id取得文件路径
	_, diskdirectory, err := GetUrlPath(product.ProjectId)
	if err != nil {
		logs.Error(err)
		c.Data["json"] = map[string]interface{}{"info": "ERROR", "state": "ERROR", "data": "获取项目路径错误！"}
		c.ServeJSON()
		return
	}

	//获取上传的文件
	_, h, err := c.GetFile("file")
	if err != nil {
		logs.Error(err)
		c.Data["json"] = map[string]interface{}{"info": "ERROR", "state": "ERROR", "data": "上传文件错误！"}
		c.ServeJSON()
		return
	}
	if h != nil {
		file_path := filepath.Clean(h.Filename)
		clean_file_path := strings.TrimPrefix(filepath.Join(string(filepath.Separator), file_path), string(filepath.Separator))
		// file_path = DiskDirectory + "/" + clean_file_path
		//保存附件
		file_path = diskdirectory + "/" + clean_file_path // attachment.FileName
		// f.Close() // 关闭上传的文件，不然的话会出现临时文件不能清除的情况
		//存入文件夹
		err = c.SaveToFile("file", file_path) //.Join("attachment", attachment)) //存文件    WaterMark(filepath)    //给文件加水印
		if err != nil {
			logs.Error(err)
		} else {
			c.Data["json"] = map[string]interface{}{"state": "SUCCESS", "title": attachment, "original": attachment}
			c.ServeJSON()
		}
	}
}

// 向某个侧栏id下新建dwg文件
func (c *AttachController) NewDwg() {
	uname, _, uid, _, _ := checkprodRole(c.Ctx)
	user, err := models.GetUserByUsername(uname)
	if err != nil {
		logs.Error(err)
	}
	var catalog models.PostMerit
	var news string
	var cid int64

	pid := c.GetString("pid")
	code := c.GetString("code")
	title := c.GetString("title")
	// subtext := c.GetString("subtext")
	label := c.GetString("label")
	principal := c.GetString("principal")
	// relevancy := c.GetString("relevancy")
	// content := c.GetString("content")

	//id转成64为
	pidNum, err := strconv.ParseInt(pid, 10, 64)
	if err != nil {
		logs.Error(err)
	}
	//根据pid查出项目id
	proj, err := models.GetProj(pidNum)
	if err != nil {
		logs.Error(err)
	}
	var topprojectid int64
	if proj.ParentIdPath != "" {
		parentidpath := strings.Replace(strings.Replace(proj.ParentIdPath, "#$", "-", -1), "$", "", -1)
		parentidpath1 := strings.Replace(parentidpath, "#", "", -1)
		patharray := strings.Split(parentidpath1, "-")
		topprojectid, err = strconv.ParseInt(patharray[0], 10, 64)
		if err != nil {
			logs.Error(err)
		}
	} else {
		topprojectid = proj.Id
	}
	//根据项目id添加成果code, title, label, principal, content string, projectid int64
	Id, err := models.AddProduct(code, title, label, principal, uid, pidNum, topprojectid)
	if err != nil {
		logs.Error(err)
	}

	//*****添加成果关联信息
	// if relevancy != "" {
	// 	_, err = models.AddRelevancy(Id, relevancy)
	// 	if err != nil {
	// 		logs.Error(err)
	// 	}
	// }
	//*****添加成果关联信息结束

	//成果写入postmerit表，准备提交merit*********
	Number, Name, DesignStage, Section, err := GetProjTitleNumber(pidNum)
	if err != nil {
		logs.Error(err)
	}
	catalog.ProjectNumber = Number
	catalog.ProjectName = Name
	catalog.DesignStage = DesignStage
	catalog.Section = Section

	catalog.Tnumber = code
	catalog.Name = title
	catalog.Count = 1
	catalog.Drawn = user.Nickname
	catalog.Designd = user.Nickname
	catalog.Author = uname
	catalog.Drawnratio = 0.4
	catalog.Designdratio = 0.4

	const lll = "2006-01-02"
	convdate := time.Now().Format(lll)
	t1, err := time.Parse(lll, convdate) //这里t1要是用t1:=就不是前面那个t1了
	if err != nil {
		logs.Error(err)
	}
	catalog.Datestring = convdate
	catalog.Date = t1

	catalog.Created = time.Now() //.Add(+time.Duration(hours) * time.Hour)
	catalog.Updated = time.Now() //.Add(+time.Duration(hours) * time.Hour)

	catalog.Complex = 1
	catalog.State = 0
	//生成提交merit的清单结束*******************

	//将dwg文件添加到成果id下
	//如果附件名称相同，则覆盖上传，但数据库不追加
	attachmentname := code + title + ".dwg"
	aid, err := models.AddAttachment(attachmentname, 0, 0, Id)
	if err != nil {
		logs.Error(err)
	} else {
		//存入文件夹
		// 	err = c.SaveToFile("file", filepath) //.Join("attachment", attachment)) //存文件    WaterMark(filepath)    //给文件加水印
		// 	if err != nil {
		// 		logs.Error(err)
		// 	}
		// 	c.Data["json"] = map[string]interface{}{"state": "SUCCESS", "title": attachment, "original": attachment, "url": Url + "/" + attachment}
		// 	c.ServeJSON()
		// }
		// aid, err := models.AddArticle(subtext, content, Id)
		// if err != nil {
		// 	logs.Error(err)
		// } else {
		//生成提交merit的清单*******************
		cid, err, news = models.AddPostMerit(catalog)
		if err != nil {
			logs.Error(err)
		} else {
			link1 := "/downloadattachment?id=" + strconv.FormatInt(aid, 10) //附件链接地址
			_, err = models.AddCatalogLink(cid, link1)
			if err != nil {
				logs.Error(err)
			}
			data := news
			c.Ctx.WriteString(data)
		}
		//生成提交merit的清单结束*******************
		c.Data["json"] = map[string]interface{}{"state": "SUCCESS", "Id": aid}
		c.ServeJSON()
	}
}

// 向某个侧栏id下添加成果——用于第二种添加，多附件模式
func (c *AttachController) AddAttachment2() {
	_, _, uid, isadmin, isLogin := checkprodRole(c.Ctx)
	if !isLogin {
		// route := c.Ctx.Request.URL.String()
		// c.Data["Url"] = route
		// c.Redirect("/roleerr?url="+route, 302)
		c.Data["json"] = "未登陆"
		c.ServeJSON()
		return
	}
	useridstring := strconv.FormatInt(uid, 10)
	pid := c.GetString("pid")
	//id转成64位
	idNum, err := strconv.ParseInt(pid, 10, 64)
	if err != nil {
		logs.Error(err)
	}
	//取项目本身
	category, err := models.GetProj(idNum)
	if err != nil {
		logs.Error(err)
	}

	var topprojectid int64
	if category.ParentId != 0 { //如果不是根目录
		parentidpath := strings.Replace(strings.Replace(category.ParentIdPath, "#$", "-", -1), "$", "", -1)
		parentidpath1 := strings.Replace(parentidpath, "#", "", -1)
		patharray := strings.Split(parentidpath1, "-")
		topprojectid, err = strconv.ParseInt(patharray[0], 10, 64)
		if err != nil {
			logs.Error(err)
		}
	} else {
		topprojectid = category.Id
	}

	projectuser, err := models.GetProjectUser(topprojectid)
	if err != nil {
		logs.Error(err)
	}
	var PostPermission bool
	if projectuser.Id == uid || isadmin {
		PostPermission = true
	} else {
		//2.取得侧栏目录路径——路由id
		//2.1 根据id取得路由
		var projurls string
		proj, err := models.GetProj(idNum)
		if err != nil {
			logs.Error(err)
		}
		if proj.ParentId == 0 { //如果是项目根目录
			projurls = "/" + strconv.FormatInt(proj.Id, 10)
		} else {
			// projurls = "/" + strings.Replace(proj.ParentIdPath, "-", "/", -1) + "/" + strconv.FormatInt(proj.Id, 10)
			projurls = "/" + strings.Replace(strings.Replace(proj.ParentIdPath, "#", "/", -1), "$", "", -1) + strconv.FormatInt(proj.Id, 10)
		}

		if res, _ := e.Enforce(useridstring, projurls+"/", "POST", ".1"); res {
			PostPermission = true
		}
	}
	if PostPermission {

		// meritbasic, err := models.GetMeritBasic()
		// if err != nil {
		// 	logs.Error(err)
		// }
		// var catalog models.PostMerit
		// var news string
		// var cid int64

		//解析表单
		// pid := c.GetString("pid")
		// beego.Info(pid)
		//pid转成64为
		// pidNum, err := strconv.ParseInt(pid, 10, 64)
		// if err != nil {
		// 	logs.Error(err)
		// }
		//根据pid查出项目id
		// proj, err := models.GetProj(pidNum)
		// if err != nil {
		// 	logs.Error(err)
		// }
		// var topprojectid int64
		// if proj.ParentIdPath != "" {
		// 	parentidpath := strings.Replace(strings.Replace(proj.ParentIdPath, "#$", "-", -1), "$", "", -1)
		// 	parentidpath1 := strings.Replace(parentidpath, "#", "", -1)
		// 	patharray := strings.Split(parentidpath1, "-")
		// 	topprojectid, err = strconv.ParseInt(patharray[0], 10, 64)
		// 	if err != nil {
		// 		logs.Error(err)
		// 	}
		// } else {
		// 	topprojectid = proj.Id
		// }
		prodcode := c.GetString("prodcode") //和上面那个区别仅仅在此处而已
		prodname := c.GetString("prodname")
		prodlabel := c.GetString("prodlabel")
		prodprincipal := c.GetString("prodprincipal")
		relevancy := c.GetString("relevancy")
		size := c.GetString("size")
		filesize, err := strconv.ParseInt(size, 10, 64)
		if err != nil {
			logs.Error(err)
		}
		filesize = filesize / 1000.0
		// proj, err := models.GetProj(pidNum)
		// if err != nil {
		// 	logs.Error(err)
		// }
		//根据proj的Id
		Url, DiskDirectory, err := GetUrlPath(idNum)
		if err != nil {
			logs.Error(err)
		}
		// Number, Name, DesignStage, Section, err := GetProjTitleNumber(idNum)
		// if err != nil {
		// 	logs.Error(err)
		// }
		// catalog.ProjectNumber = Number
		// catalog.ProjectName = Name
		// catalog.DesignStage = DesignStage
		// catalog.Section = Section

		//获取上传的文件
		_, h, err := c.GetFile("file")
		if err != nil {
			logs.Error(err)
		}
		// var filesize int64
		if h != nil {

			//保存附件
			// attachment = h.Filename
			// f.Close()// 关闭上传的文件，不然的话会出现临时文件不能清除的情况
			//存入成果数据库
			//如果编号重复，则不写入，值返回Id值。
			//根据id添加成果code, title, label, principal, content string, projectid int64
			prodId, err := models.AddProduct(prodcode, prodname, prodlabel, prodprincipal, uid, idNum, topprojectid)
			if err != nil {
				logs.Error(err)
			}
			//20190728改名，替换文件名中的#和斜杠
			filename2 := strings.Replace(h.Filename, "#", "号", -1)
			filename2 = strings.Replace(filename2, "/", "-", -1)
			// FileSuffix := path.Ext(h.Filename)
			// attachmentname := filename2 + FileSuffix
			file_name := filepath.Clean(filename2)
			clean_file_name := strings.TrimPrefix(filepath.Join(string(filepath.Separator), file_name), string(filepath.Separator))

			file_path := DiskDirectory + "/" + clean_file_name // filename2
			//*****添加成果关联信息
			if relevancy != "" {
				array := strings.Split(relevancy, ",")
				for _, v := range array {
					_, err = models.AddRelevancy(prodId, v)
					if err != nil {
						logs.Error(err)
					}
				}
			}
			//*****添加成果关联信息结束

			//20200918注释掉以下部分。成果写入postmerit表，准备提交merit*********
			// catalog.Tnumber = prodcode
			// catalog.Name = prodname
			// catalog.Count = 1
			// catalog.Drawn = meritbasic.Nickname
			// catalog.Designd = meritbasic.Nickname
			// catalog.Author = meritbasic.Username
			// catalog.Drawnratio = 0.4
			// catalog.Designdratio = 0.4
			// const lll = "2006-01-02"
			// convdate := time.Now().Format(lll)
			// t1, err := time.Parse(lll, convdate) //这里t1要是用t1:=就不是前面那个t1了
			// if err != nil {
			// 	logs.Error(err)
			// }
			// catalog.Datestring = convdate
			// catalog.Date = t1
			// catalog.Created = time.Now() //.Add(+time.Duration(hours) * time.Hour)
			// catalog.Updated = time.Now() //.Add(+time.Duration(hours) * time.Hour)
			// catalog.Complex = 1
			// catalog.State = 0
			// cid, err, news = models.AddPostMerit(catalog)
			// if err != nil {
			// 	logs.Error(err)
			// } else {
			// 	link1 := Url + "/" + filename2 //附件链接地址
			// 	_, err = models.AddCatalogLink(cid, link1)
			// 	if err != nil {
			// 		logs.Error(err)
			// 	}
			// 	data := news
			// 	c.Ctx.WriteString(data)
			// }
			//生成提交merit的清单结束*******************

			//把成果id作为附件的parentid，把附件的名称等信息存入附件数据库
			//如果附件名称相同，则覆盖上传，但数据库不追加
			_, err = models.AddAttachment(filename2, filesize, 0, prodId)
			if err != nil {
				logs.Error(err)
				c.Data["json"] = map[string]interface{}{"state": "WRONG存入数据库出错", "title": err, "original": "", "url": ""}
				c.ServeJSON()
			} else {
				//建立目录，并返回作为父级目录
				err = os.MkdirAll(DiskDirectory+"/", 0777) //..代表本当前exe文件目录的上级，.表示当前目录，没有.表示盘的根目录
				if err != nil {
					logs.Error(err)
				}
				err = c.SaveToFile("file", file_path) //.Join("attachment", attachment)) //存文件    WaterMark(path)    //给文件加水印
				if err != nil {
					logs.Error(err)
					c.Data["json"] = map[string]interface{}{"state": "WRONG保存文件出错", "title": err, "original": "", "url": ""}
					c.ServeJSON()
				} else {
					c.Data["json"] = map[string]interface{}{"state": "SUCCESS", "title": filename2, "original": filename2, "url": Url + "/" + filename2}
					c.ServeJSON()
				}
			}
		}
	} else {
		c.Data["json"] = map[string]interface{}{"state": "ERROR", "title": "非管理员、非本人、未赋予添加权限"}
		c.ServeJSON()
	}
}

// 向一个成果id下追加附件
func (c *AttachController) UpdateAttachment() {
	// _, role := checkprodRole(c.Ctx)
	// if role == 1 {
	//解析表单
	pid := c.GetString("pid")
	size := c.GetString("size")
	filesize, err := strconv.ParseInt(size, 10, 64)
	if err != nil {
		logs.Error(err)
	}
	filesize = filesize / 1000.0
	// beego.Info(pid)
	//pid转成64为
	pidNum, err := strconv.ParseInt(pid, 10, 64)
	if err != nil {
		logs.Error(err)
	}

	prod, err := models.GetProd(pidNum)
	if err != nil {
		logs.Error(err)
	}
	//根据proj的Id
	Url, DiskDirectory, err := GetUrlPath(prod.ProjectId)
	if err != nil {
		logs.Error(err)
	}
	//获取上传的文件
	_, h, err := c.GetFile("file")
	if err != nil {
		logs.Error(err)
	}

	if h != nil {
		file_name := filepath.Clean(h.Filename)
		clean_file_name := strings.TrimPrefix(filepath.Join(string(filepath.Separator), file_name), string(filepath.Separator))
		//保存附件
		file_path := DiskDirectory + "/" + clean_file_name // h.Filename // 关闭上传的文件，不然的话会出现临时文件不能清除的情况
		// filesize, _ = FileSize(path)
		// filesize = filesize / 1000.0
		//把成果id作为附件的parentid，把附件的名称等信息存入附件数据库
		//如果附件名称相同，则覆盖上传，但数据库不追加
		_, err = models.AddAttachment(clean_file_name, filesize, 0, pidNum)
		if err != nil {
			logs.Error(err)
		} else {
			err = c.SaveToFile("file", file_path) //.Join("attachment", attachment)) //存文件    WaterMark(path)    //给文件加水印
			if err != nil {
				logs.Error(err)
			}
			c.Data["json"] = map[string]interface{}{"state": "SUCCESS", "title": clean_file_name, "original": clean_file_name, "url": Url + "/" + clean_file_name}
			c.ServeJSON()
		}
	}
	// } else {
	// 	route := c.Ctx.Request.URL.String()
	// 	c.Data["Url"] = route
	// 	c.Redirect("/roleerr?url="+route, 302)
	// 	// c.Redirect("/roleerr", 302)
	// 	return
	// }
}

// 删除附件——这个用于针对删除一个附件
func (c *AttachController) DeleteAttachment() {
	_, _, uid, isadmin, isLogin := checkprodRole(c.Ctx)
	if !isLogin {
		// route := c.Ctx.Request.URL.String()
		// c.Data["Url"] = route
		// c.Redirect("/roleerr?url="+route, 302)
		c.Data["json"] = "未登陆"
		c.ServeJSON()
		return
	}
	useridstring := strconv.FormatInt(uid, 10)
	// pid := c.GetString("pid")
	ids := c.GetString("ids")
	array := strings.Split(ids, ",")
	//id转成64位
	idNum, err := strconv.ParseInt(array[0], 10, 64)
	if err != nil {
		logs.Error(err)
	}
	//取项目本身
	category, err := models.GetProj(idNum)
	if err != nil {
		logs.Error(err)
	}

	var topprojectid int64
	if category.ParentId != 0 { //如果不是根目录
		parentidpath := strings.Replace(strings.Replace(category.ParentIdPath, "#$", "-", -1), "$", "", -1)
		parentidpath1 := strings.Replace(parentidpath, "#", "", -1)
		patharray := strings.Split(parentidpath1, "-")
		topprojectid, err = strconv.ParseInt(patharray[0], 10, 64)
		if err != nil {
			logs.Error(err)
		}
	} else {
		topprojectid = category.Id
	}

	projectuser, err := models.GetProjectUser(topprojectid)
	if err != nil {
		logs.Error(err)
	}
	var DeletePermission bool
	if projectuser.Id == uid || isadmin {
		DeletePermission = true
	} else {
		//2.取得侧栏目录路径——路由id
		//2.1 根据id取得路由
		var projurls string
		proj, err := models.GetProj(idNum)
		if err != nil {
			logs.Error(err)
		}
		if proj.ParentId == 0 { //如果是项目根目录
			projurls = "/" + strconv.FormatInt(proj.Id, 10)
		} else {
			// projurls = "/" + strings.Replace(proj.ParentIdPath, "-", "/", -1) + "/" + strconv.FormatInt(proj.Id, 10)
			projurls = "/" + strings.Replace(strings.Replace(proj.ParentIdPath, "#", "/", -1), "$", "", -1) + strconv.FormatInt(proj.Id, 10)
		}

		if res, _ := e.Enforce(useridstring, projurls+"/", "DELETE", ".1"); res {
			DeletePermission = true
		}
	}
	if DeletePermission {

		for _, v := range array {
			// pid = strconv.FormatInt(v1, 10)
			//id转成64位
			idNum, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				logs.Error(err)
			}
			//取得附件的成果id——再取得成果的项目目录id——再取得路径
			attach, err := models.GetAttachbyId(idNum)
			if err != nil {
				logs.Error(err)
			}
			prod, err := models.GetProd(attach.ProductId)
			if err != nil {
				logs.Error(err)
			}
			//根据proj的id
			_, DiskDirectory, err := GetUrlPath(prod.ProjectId)
			if err != nil {
				logs.Error(err)
			} else if DiskDirectory != "" {
				path := DiskDirectory + "/" + attach.FileName
				err = os.Remove(path)
				if err != nil {
					logs.Error(err)
				}
				err = models.DeleteAttachment(idNum)
				if err != nil {
					logs.Error(err)
				}
			}
		}
		c.Data["json"] = "ok"
		c.ServeJSON()
	} else {
		c.Data["json"] = "非管理员、非本人、未赋予删除权限"
		c.ServeJSON()
	}
	// } else {
	// 	route := c.Ctx.Request.URL.String()
	// 	c.Data["Url"] = route
	// 	c.Redirect("/roleerr?url="+route, 302)
	// 	// c.Redirect("/roleerr", 302)
	// 	return
	// }
}

// 下载附件——这个仅是测试
// func (c *AttachController) ImageFilter() {
func ImageFilter(ctx *context.Context) {
	// token := path.Base(ctx.Request.RequestURI)
	// split token and file ext
	// var filePath string
	// if i := strings.IndexRune(token, '.'); i == -1 {
	// 	return
	// } else {
	// 	filePath = token[i+1:]
	// 	token = token[:i]
	// }
	// decode token to file path
	// var image models.Image
	// if err := image.DecodeToken(token); err != nil {
	// 	beego.Info(err)
	// 	return
	// }
	// file real path
	// filePath = GetUrlPath(&image) + filePath
	//1.url处理中文字符路径，[1:]截掉路径前面的/斜杠
	filePath := path.Base(ctx.Request.RequestURI)
	// filePath, err := url.QueryUnescape(c.Ctx.Request.RequestURI[1:]) //  attachment/SL2016测试添加成果/A/FB/1/Your First Meteor Application.pdf
	// if err != nil {
	// 	logs.Error(err)
	// }
	// if x-send on then set header and http status
	// fall back use proxy serve file
	// if setting.ImageXSend {
	// 	ext := filepath.Ext(filePath)
	// 	ctx.Output.ContentType(ext)
	// 	ctx.Output.Header(setting.ImageXSendHeader, "/"+filePath)
	// 	ctx.Output.SetStatus(200)
	// } else {
	// direct serve file use go
	http.ServeFile(ctx.ResponseWriter, ctx.Request, filePath)
	// http.ServeFile(c.Ctx.ResponseWriter, c.Ctx.Request, filePath)
	// }
}

// 根据权限查看附件/downloadattachment?id=
func (c *AttachController) DownloadAttachment() {
	username, _, uid, isadmin, isLogin := checkprodRole(c.Ctx)
	if !isLogin {
		// route := c.Ctx.Request.URL.String()
		// c.Data["Url"] = route
		// c.Redirect("/roleerr?url="+route, 302)
		c.Data["json"] = "未登陆"
		c.ServeJSON()
		return
	}
	useridstring := strconv.FormatInt(uid, 10)
	//4.取得客户端用户名
	var projurl string
	// var usersessionid string
	// if isLogin {
	// 	usersessionid = c.Ctx.Input.Cookie("hotqinsessionid")
	// }
	id := c.GetString("id")
	//pid转成64为
	idNum, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		logs.Error(err)
		utils.FileLogs.Error(c.Ctx.Input.IP() + " 转换id " + err.Error())
	}

	//根据附件id取得附件的prodid，路径
	attachment, err := models.GetAttachbyId(idNum)
	if err != nil {
		logs.Error(err)
		utils.FileLogs.Error(c.Ctx.Input.IP() + " 查询附件 " + err.Error())
	}

	product, err := models.GetProd(attachment.ProductId)
	if err != nil {
		logs.Error(err)
		utils.FileLogs.Error(c.Ctx.Input.IP() + " 用附件查询成果 " + err.Error())
	}

	//取项目本身
	category, err := models.GetProj(product.ProjectId)
	if err != nil {
		logs.Error(err)
	}

	var topprojectid int64
	if category.ParentId != 0 { //如果不是根目录
		parentidpath := strings.Replace(strings.Replace(category.ParentIdPath, "#$", "-", -1), "$", "", -1)
		parentidpath1 := strings.Replace(parentidpath, "#", "", -1)
		patharray := strings.Split(parentidpath1, "-")
		topprojectid, err = strconv.ParseInt(patharray[0], 10, 64)
		if err != nil {
			logs.Error(err)
		}
	} else {
		topprojectid = category.Id
	}

	projectuser, err := models.GetProjectUser(topprojectid)
	if err != nil {
		logs.Error(err)
	}
	var isme bool
	if projectuser.Id == uid {
		isme = true
	}

	//根据projid取出路径
	// proj, err := models.GetProj(product.ProjectId)
	// if err != nil {
	// 	logs.Error(err)
	// 	utils.FileLogs.Error(err.Error())
	// }
	if category.ParentIdPath == "" || category.ParentIdPath == "$#" {
		projurl = "/" + strconv.FormatInt(category.Id, 10) + "/"
	} else {
		// projurl = "/" + strings.Replace(proj.ParentIdPath, "-", "/", -1) + "/" + strconv.FormatInt(proj.Id, 10) + "/"
		projurl = "/" + strings.Replace(strings.Replace(category.ParentIdPath, "#", "/", -1), "$", "", -1) + strconv.FormatInt(category.Id, 10) + "/"
	}
	//由proj id取得url
	fileurl, _, err := GetUrlPath(product.ProjectId)
	if err != nil {
		logs.Error(err)
		utils.FileLogs.Error(c.Ctx.Input.IP() + " 查询成果路径 " + err.Error())
	}

	fileext := path.Ext(attachment.FileName)
	matched, err := regexp.MatchString("\\.*[m|M][c|C][d|D]", fileext)
	if err != nil {
		logs.Error(err)
	}
	// beego.Info(matched)
	if matched {
		c.Data["json"] = map[string]interface{}{"info": "ERROR", "data": "不能下载mcd文件!"}
		c.ServeJSON()
		return
	}
	matched, err = regexp.MatchString("\\.*[f|F][c|C][s|S][t|T][d|D]", fileext)
	if err != nil {
		logs.Error(err)
	}
	// beego.Info(matched)
	if matched {
		c.Data["json"] = map[string]interface{}{"info": "ERROR", "data": "不能下载fcstd文件!"}
		c.ServeJSON()
		return
	}
	switch fileext {
	case ".JPG", ".jpg", ".jpeg", ".JPEG", ".png", ".PNG", ".bmp", ".BMP":
		// c.Ctx.Output.Download(fileurl + "/" + attachment.FileName)
		http.ServeFile(c.Ctx.ResponseWriter, c.Ctx.Request, fileurl+"/"+attachment.FileName)
	case ".mp4", ".MP4":
		res, err := e.Enforce(useridstring, projurl, c.Ctx.Request.Method, fileext)
		if err != nil {
			logs.Error(err)
		}
		if res || isadmin || isme { //+ strconv.Itoa(c.Ctx.Input.Port())
			mp4link := "/" + fileurl + "/" + attachment.FileName
			if err != nil {
				logs.Error(err)
				utils.FileLogs.Error(c.Ctx.Input.IP() + " 获取mp4路径 " + err.Error())
			}
			c.Data["FileName"] = attachment.FileName
			c.Data["Id"] = id
			c.Data["Mp4Link"] = mp4link
			// c.Data["Sessionid"] = usersessionid
			c.TplName = "flv.tpl"
		} else {
			route := c.Ctx.Request.URL.String()
			c.Data["Url"] = route
			c.Redirect("/roleerr?url="+route, 302)
			return
		}
	// case ".dwg", ".DWG": //保留，dwg在线阅览模式！！！！
	// 	if e.Enforce(useridstring, projurl, c.Ctx.Request.Method, fileext) || isadmin || isme { //+ strconv.Itoa(c.Ctx.Input.Port())
	// 		dwglink, err := url.ParseRequestURI(c.Ctx.Input.Scheme() + "://" + c.Ctx.Input.IP() + ":" + strconv.Itoa(c.Ctx.Input.Port()) + "/" + fileurl + "/" + attachment.FileName)
	// 		if err != nil {
	// 			logs.Error(err)
	// 			utils.FileLogs.Error(c.Ctx.Input.IP() + " 获取dwg路径 " + err.Error())
	// 		}
	// 		c.Data["FileName"] = attachment.FileName
	// 		c.Data["Id"] = id
	// 		c.Data["DwgLink"] = dwglink
	// 		beego.Info(dwglink)
	// 		c.Data["Sessionid"] = usersessionid
	// 		c.TplName = "dwg.tpl"
	// 	} else {
	// 		route := c.Ctx.Request.URL.String()
	// 		c.Data["Url"] = route
	// 		c.Redirect("/roleerr?url="+route, 302)
	// 		return
	// 	}
	default:
		res2, err := e.Enforce(useridstring, projurl, c.Ctx.Request.Method, fileext)
		if err != nil {
			logs.Error(err)
		}
		if res2 || isadmin || isme {
			// http.ServeFile(c.Ctx.ResponseWriter, c.Ctx.Request, filePath)//这样写下载的文件名称不对
			// c.Redirect(url+"/"+attachment.FileName, 302)
			c.Ctx.Output.Download(fileurl + "/" + attachment.FileName)
			logs.Info("下载……" + fileurl + "/" + attachment.FileName)
			utils.FileLogs.Info(username + " " + "download" + " " + fileurl + "/" + attachment.FileName)
		} else {
			utils.FileLogs.Info(c.Ctx.Input.IP() + "want " + "download" + " " + fileurl + "/" + attachment.FileName)
			route := c.Ctx.Request.URL.String()
			c.Data["Url"] = route
			c.Redirect("/roleerr?url="+route, 302)
			// c.Redirect("/roleerr", 302)
			return
		}
	}
	// config := make(map[string]interface{})
	// config["filename"] = "e:/golang/go_pro/logs/logcollect.log"
	// config["filename"] = "e://golang//go_pro//logs//logcollect.log"
	// config["filename"] = "e:/golang/go_pro/logs/logcollect.log"
	// config["level"] = logs.LevelDebug
	// configStr, err := json.Marshal(config)
	// if err != nil {
	// 	fmt.Println("marshal failed, err:", err)
	// 	return
	// }
	// logs.SetLogger(logs.AdapterFile, string(configStr))
	// utils.FileLogs.Warn("this is a warn, my name is %s", map[string]int{"key": 2016})
	// utils.FileLogs.Critical("oh,crash")
	// logs.Close()
	// 日志有以下7个级别                  对应的方法
	// LevelEmergency = 0      --> logs.Emergency()
	// LevelAlert = 1          --> logs.Alert()
	// LevelCritical = 2       --> logs.Critical()
	// LevelError = 3          --> logs.Error()
	// LevelWarning = 4        --> logs.Warning()
	// LevelNotice = 5         --> logs.Notice()
	// LevelInformational = 6  --> logs.Informational()
	// LevelDebug = 7          --> logs.Debug()
	// utils.FileLogs.Info("this is a file log with info.")
	// utils.FileLogs.Debug("this is a file log with debug.")
	// utils.FileLogs.Alert("this is a file log with alert.")
	// utils.FileLogs.Error("this is a file log with error.")
	// utils.FileLogs.Trace("this is a file log with trace.")
}

// 目前有文章中的图片、成果中文档的预览、onlyoffice中的文档协作、pdf中的附件路径等均采用绝对路径型式
// 文章中的附件呢？
// default中的pdf页面中的{{.pdflink}}，绝对路径
//
//	type Session struct {
//		Session int
//	}
//
// attachment/路径/附件名称
func (c *AttachController) Attachment() {
	//如果url带了sessionid,就能取到uid等信息
	var useridstring string
	_, _, uid, isadmin, islogin := checkprodRole(c.Ctx)
	// logs.Info(uid)
	// logs.Info(islogin)
	useridstring = strconv.FormatInt(uid, 10)

	//1.url处理中文字符路径，[1:]截掉路径前面的/斜杠
	// filePath := path.Base(c.Ctx.Request.RequestURI)
	filePath, err := url.QueryUnescape(c.Ctx.Request.RequestURI[1:]) //attachment/SL2016测试添加成果/A/FB/1/Your First Meteor Application.pdf
	if err != nil {
		logs.Error(err)
	}
	// beego.Info(filePath)
	//attachment/standard/SL/SLZ 5077-2016水工建筑物荷载设计规范.pdf
	if strings.Contains(filePath, "?") { //hotqinsessionid=
		filePathtemp := strings.Split(filePath, "?")
		filePath = filePathtemp[0]
	}
	fileext := path.Ext(filePath)
	filepath1 := path.Dir(filePath)
	array := strings.Split(filepath1, "/")
	// logs.Info(array[1])
	// beego.Info(fileext)
	matched, err := regexp.MatchString("\\.*[m|M][c|C][d|D]", fileext)
	if err != nil {
		logs.Error(err)
	}
	// logs.Info(matched)
	if matched {
		c.Data["json"] = map[string]interface{}{"info": "ERROR", "data": "不能下载mcd文件!"}
		c.ServeJSON()
		return
	}
	matched, err = regexp.MatchString("\\.*[f|F][c|C][s|S][t|T][d|D]", fileext)
	if err != nil {
		logs.Error(err)
	}
	// logs.Info(matched)
	if matched {
		c.Data["json"] = map[string]interface{}{"info": "ERROR", "data": "不能下载fcstd文件!"}
		c.ServeJSON()
		return
	}
	// logs.Info(matched)
	// 这里设计不合理，只要是注册人员均可查阅计算稿历史，不好。
	if array[1] == "standard" || (array[1] == "math" && fileext == ".pdf") || (array[1] == "pass_excel" && fileext == ".pdf") || (array[1] == "pass_ansys" && fileext == ".pdf") {
		if !islogin {
			// beego.Info(!islogin)
			route := c.Ctx.Request.URL.String()
			c.Data["Url"] = route
			c.Redirect("/roleerr?url="+route, 302)
			return
		} else {
			http.ServeFile(c.Ctx.ResponseWriter, c.Ctx.Request, filePath)
			return
		}
	}
	//查出所有项目
	var pid int64
	proj, err := models.GetProjects()
	// 循环，项目编号+名称=array[1]则其id
	// 20220102bug_如果项目编号和名称相同，则出错！！！！！！
	for _, v := range proj {
		if v.Code+v.Title == array[1] {
			pid = v.Id
			break
		}
	}

	//取项目本身
	category, err := models.GetProj(pid)
	if err != nil {
		logs.Error(err)
	}

	var topprojectid int64
	if category.ParentId != 0 { //如果不是根目录
		parentidpath := strings.Replace(strings.Replace(category.ParentIdPath, "#$", "-", -1), "$", "", -1)
		parentidpath1 := strings.Replace(parentidpath, "#", "", -1)
		patharray := strings.Split(parentidpath1, "-")
		topprojectid, err = strconv.ParseInt(patharray[0], 10, 64)
		if err != nil {
			logs.Error(err)
		}
	} else {
		topprojectid = category.Id
	}

	projectuser, err := models.GetProjectUser(topprojectid)
	if err != nil {
		logs.Error(err)
	}
	// beego.Info(projectuser.Id)
	var isme bool
	if projectuser.Id == uid && uid != 0 {
		isme = true
	}

	var projurl models.Project
	projurls := "/" + strconv.FormatInt(pid, 10)
	//id作为parentid+array[2]
	if len(array) > 1 {
		for i := 2; i < len(array); i++ {
			if i == 2 {
				// beego.Info(pid)
				projurl, err = models.GetProjbyParentidTitle(pid, array[i])
				if err != nil {
					logs.Error(err)
				}
			} else {
				projurl, err = models.GetProjbyParentidTitle(projurl.Id, array[i])
				if err != nil {
					// logs.Error(err)
				}
			}
			projurls = projurls + "/" + strconv.FormatInt(projurl.Id, 10)
		}
	}

	switch fileext {
	case ".JPG", ".jpg", ".jpeg", ".JPEG", ".png", ".PNG", ".bmp", ".BMP", ".mp4", ".MP4", ".glb", ".gltf", ".bin", ".svg":
		http.ServeFile(c.Ctx.ResponseWriter, c.Ctx.Request, filePath)
		// case ".dwg", ".DWG":
		// http.ServeFile(c.Ctx.ResponseWriter, c.Ctx.Request, filePath)
		// case ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx":
		// beego.Info("ok")
		// http.ServeFile(c.Ctx.ResponseWriter, c.Ctx.Request, filePath)
		// logs.Info(filePath)
		// logs.Info(useridstring)
	//这里缺少权限设置！！！！！！！！！！！
	default:
		res, err := e.Enforce(useridstring, projurls+"/", c.Ctx.Request.Method, fileext)
		if err != nil {
			logs.Error(err)
		}
		// logs.Info(res)
		if res || isadmin || isme {
			// beego.Info(e.Enforce(useridstring, projurls+"/", c.Ctx.Request.Method, fileext))
			http.ServeFile(c.Ctx.ResponseWriter, c.Ctx.Request, filePath) //这样写下载的文件名称不对
			// beego.Info(isadmin)
			// c.Redirect(url+"/"+attachment.FileName, 302)
			// c.Ctx.Output.Download(fileurl + "/" + attachment.FileName)
		} else {
			// beego.Info(useridstring)
			route := c.Ctx.Request.URL.String()
			c.Data["Url"] = route
			c.Redirect("/roleerr?url="+route, 302)
			// c.Redirect("/roleerr", 302)
			return
		}
	}
	// switch fileext {
	// case ".pdf", ".PDF", ".JPG", ".jpg", ".png", ".PNG", ".bmp", ".BMP", ".mp4", ".MP4":
	// 	if role < 5 {
	// 		http.ServeFile(c.Ctx.ResponseWriter, c.Ctx.Request, filePath)
	// 	} else {
	// 		route := c.Ctx.Request.URL.String()
	// 		c.Data["Url"] = route
	// 		c.Redirect("/roleerr?url="+route, 302)
	// 		return
	// 	}
	// default:
	// 	if role <= 2 {
	// 		http.ServeFile(c.Ctx.ResponseWriter, c.Ctx.Request, filePath)
	// 	} else {
	// 		route := c.Ctx.Request.URL.String()
	// 		c.Data["Url"] = route
	// 		c.Redirect("/roleerr?url="+route, 302)
	// 		return
	// 	}
	// }
}

// 首页轮播图片给予任何权限
func (c *AttachController) GetCarousel() {
	//1.url处理中文字符路径，[1:]截掉路径前面的/斜杠
	// filePath := path.Base(ctx.Request.RequestURI)
	filePath, err := url.QueryUnescape(c.Ctx.Request.RequestURI[1:]) //  attachment/SL2016测试添加成果/A/FB/1/Your First Meteor Application.pdf
	if err != nil {
		logs.Error(err)
	}

	fileext := path.Ext(filePath)
	matched, err := regexp.MatchString("\\.*[m|M][c|C][d|D]", fileext)
	if err != nil {
		logs.Error(err)
	}
	// beego.Info(matched)
	if matched {
		c.Data["json"] = map[string]interface{}{"info": "ERROR", "data": "不能下载mcd文件!"}
		c.ServeJSON()
		return
	}

	http.ServeFile(c.Ctx.ResponseWriter, c.Ctx.Request, filePath)
}

// 返回文件大小
func FileSize(file string) (int64, error) {
	f, e := os.Stat(file)
	if e != nil {
		return 0, e
	}
	return f.Size(), nil
}

// 根据侧栏id返回附件url和文件夹路径
func GetUrlPath(id int64) (Url, DiskDirectory string, err error) {
	var parentidpath, parentidpath1 string
	proj, err := models.GetProj(id)
	if err != nil {
		logs.Error(err)
		return "", "", err
	} else {
		//根据proj的parentIdpath
		var path string
		if proj.ParentIdPath != "" { //如果不是根目录
			// patharray := strings.Split(proj.ParentIdPath, "-")
			parentidpath = strings.Replace(strings.Replace(proj.ParentIdPath, "#$", "-", -1), "$", "", -1)
			parentidpath1 = strings.Replace(parentidpath, "#", "", -1)
			patharray := strings.Split(parentidpath1, "-")
			for _, v := range patharray {
				//pid转成64为
				idNum, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					logs.Error(err)
				}
				proj1, err := models.GetProj(idNum)
				if err != nil {
					logs.Error(err)
				} else {
					if proj1.ParentId == 0 { //如果是项目名称，则加上项目编号
						DiskDirectory = "./attachment/" + proj1.Code + proj1.Title
						Url = "attachment/" + proj1.Code + proj1.Title
					} else {
						path = proj1.Title
						DiskDirectory = DiskDirectory + "/" + path
						Url = Url + "/" + path
					}
				}
			}
			DiskDirectory = DiskDirectory + "/" + proj.Title //加上自身
			Url = Url + "/" + proj.Title
		} else { //如果是根目录
			DiskDirectory = "./attachment/" + proj.Code + proj.Title //加上自身
			Url = "attachment/" + proj.Code + proj.Title
		}
		return Url, DiskDirectory, err
	}
}

// 根据id返回项目编号，项目名称，项目阶段，项目专业
func GetProjTitleNumber(id int64) (ProjectNumber, ProjectName, DesignStage, Section string, err error) {
	var parentidpath, parentidpath1 string
	proj, err := models.GetProj(id)
	if err != nil {
		logs.Error(err)
		return "", "", "", "", err
	} else {
		//根据proj的parentIdpath
		if proj.ParentIdPath != "" { //如果不是根目录
			parentidpath = strings.Replace(strings.Replace(proj.ParentIdPath, "#$", "-", -1), "$", "", -1)
			parentidpath1 = strings.Replace(parentidpath, "#", "", -1)
			patharray := strings.Split(parentidpath1, "-")
			// patharray := strings.Split(proj.ParentIdPath, "-")
			//pid转成64位
			meritNum, err := strconv.ParseInt(patharray[0], 10, 64)
			if err != nil {
				logs.Error(err)
			}
			meritproj, err := models.GetProj(meritNum)
			if err != nil {
				logs.Error(err)
			}
			ProjectNumber = meritproj.Code
			ProjectName = meritproj.Title
			if len(patharray) > 1 {
				//pid转成64位
				meritNum1, err := strconv.ParseInt(patharray[1], 10, 64)
				if err != nil {
					logs.Error(err)
				}
				meritproj1, err := models.GetProj(meritNum1)
				if err != nil {
					logs.Error(err)
				}
				DesignStage = meritproj1.Title
			}

			if len(patharray) > 2 {
				//pid转成64位
				meritNum2, err := strconv.ParseInt(patharray[2], 10, 64)
				if err != nil {
					logs.Error(err)
				}
				meritproj2, err := models.GetProj(meritNum2)
				if err != nil {
					logs.Error(err)
				}
				Section = meritproj2.Title
			}

		} else { //如果是根目录
			ProjectNumber = proj.Code
			ProjectName = proj.Title
		}
		return ProjectNumber, ProjectName, DesignStage, Section, err
	}
}

// 网页阅览pdf文件，并带上上一页和下一页
func (c *AttachController) Pdf() {
	_, _, uid, _, _ := checkprodRole(c.Ctx)
	// logs.Info(uid)
	if uid == 0 {
		route := c.Ctx.Request.URL.String()
		c.Data["Url"] = route
		c.Redirect("/roleerr?url="+route, 302)
		return
	}

	//取得某个目录下所有pdf
	var Url string
	var pNum int64
	id := c.GetString("id")
	// beego.Info(id)
	//id转成64为
	idNum, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		logs.Error(err)
	}
	p := c.GetString("p")
	// beego.Info(p)

	var prods []*models.Product
	var projid int64
	if p == "" { //id是附件或成果id
		attach, err := models.GetAttachbyId(idNum)
		if err != nil {
			logs.Error(err)
		}
		//由成果id（后台传过来的行id）取得成果——进行排序
		prod, err := models.GetProd(attach.ProductId)
		if err != nil {
			logs.Error(err)
		}
		//由侧栏id取得所有成果和所有成果的附件pdf
		prods, err = models.GetProducts(prod.ProjectId)
		if err != nil {
			logs.Error(err)
		}

		projid = prod.ProjectId
	} else { //id是侧栏目录id
		//由侧栏id取得所有成果和所有成果的附件pdf——进行排序
		prods, err = models.GetProducts(idNum)
		if err != nil {
			logs.Error(err)
		}
		//由proj id取得url
		Url, _, err = GetUrlPath(idNum)
		if err != nil {
			logs.Error(err)
		}
		//id转成64为
		pNum, err = strconv.ParseInt(p, 10, 64)
		if err != nil {
			logs.Error(err)
		}
	}

	Attachments := make([]*models.Attachment, 0)
	for _, v := range prods {
		//根据成果id取得所有附件数量
		Attachments1, err := models.GetAttachments(v.Id) //只返回Id?没意义，全部返回
		if err != nil {
			logs.Error(err)
		}
		for _, v := range Attachments1 {
			if path.Ext(v.FileName) == ".pdf" || path.Ext(v.FileName) == ".PDF" {
				Attachments = append(Attachments, Attachments1...)
			}
		}
	}

	count := len(Attachments)
	count1 := strconv.Itoa(count)
	count2, err := strconv.ParseInt(count1, 10, 64)
	if err != nil {
		logs.Error(err)
	}
	postsPerPage := 1
	paginator := pagination.SetPaginator(c.Ctx, postsPerPage, count2)

	c.Data["paginator"] = paginator
	if p == "" {
		var p1 string
		for i, v := range Attachments {
			if v.Id == idNum {
				p1 = strconv.Itoa(i + 1)
				// PdfLink = Url + "/" + v.FileName
				break
			}
		}
		c.Redirect("/pdf?p="+p1+"&id="+strconv.FormatInt(projid, 10), 302)
	} else {
		PdfLink := Url + "/" + Attachments[pNum-1].FileName
		logs.Info(PdfLink)
		c.Data["PdfLink"] = PdfLink
		c.TplName = "web/viewer.html"
	}
	// logs := logs.NewLogger(1000)
	// logs.SetLogger("file", `{"filename":"log/test.log"}`)
	// logs.EnableFuncCallDepth(true)
	// logs.Info(c.Ctx.Input.IP() + " " + "ListAllTopic")
	// logs.Close()
	// count, _ := models.M("logoperation").Alias(`op`).Field(`count(op.id) as count`).Where(where).Count()
	// if count > 0 {
	// 	pagesize := 10
	// 	p := tools.NewPaginator(this.Ctx.Request, pagesize, count)
	// 	log, _ := models.M("logoperation").Alias(`op`).Where(where).Limit(strconv.Itoa(p.Offset()), strconv.Itoa(pagesize)).Order(`op.id desc`).Select()
	// 	this.Data["data"] = log
	// 	this.Data["paginator"] = p
	// }
}

// @Title dowload wx pdf
// @Description get wx pdf by id
// @Param id path string  true "The id of pdf"
// @Success 200 {object} models.GetAttachbyId
// @Failure 400 Invalid page supplied
// @Failure 404 pdf not found
// @router /wxpdf/:id [get]
// 有权限判断,必须取得文件所在路径，才能用casbin判断
func (c *AttachController) WxPdf() {
	// var openID string
	// openid := c.GetSession("openID")
	// beego.Info(openid)
	// if openid == nil {

	// } else {
	// 	openID = openid.(string)
	// }
	openID := c.GetSession("openID")
	// beego.Info(openID)
	if openID != nil {
		user, err := models.GetUserByOpenID(openID.(string))
		if err != nil {
			logs.Error(err)
		} else {
			useridstring := strconv.FormatInt(user.Id, 10)
			// 判断用户是否具有权限。
			id := c.Ctx.Input.Param(":id")
			//pid转成64为
			idNum, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				logs.Error(err)
			}

			//根据附件id取得附件的prodid，路径
			attachment, err := models.GetAttachbyId(idNum)
			if err != nil {
				logs.Error(err)
			}

			product, err := models.GetProd(attachment.ProductId)
			if err != nil {
				logs.Error(err)
			}
			//根据projid取出路径
			proj, err := models.GetProj(product.ProjectId)
			if err != nil {
				logs.Error(err)
				utils.FileLogs.Error(err.Error())
			}
			var projurl string
			if proj.ParentIdPath == "" || proj.ParentIdPath == "$#" {
				projurl = "/" + strconv.FormatInt(proj.Id, 10) + "/"
			} else {
				projurl = "/" + strings.Replace(strings.Replace(proj.ParentIdPath, "#", "/", -1), "$", "", -1) + strconv.FormatInt(proj.Id, 10) + "/"
			}
			//由proj id取得url
			fileurl, _, err := GetUrlPath(product.ProjectId)
			if err != nil {
				logs.Error(err)
			}
			fileext := path.Ext(attachment.FileName)
			matched, err := regexp.MatchString("\\.*[m|M][c|C][d|D]", fileext)
			if err != nil {
				logs.Error(err)
			}
			if matched {
				c.Data["json"] = map[string]interface{}{"info": "ERROR", "data": "不能下载mcd文件!"}
				c.ServeJSON()
				return
			}

			if res, _ := e.Enforce(useridstring, projurl, c.Ctx.Request.Method, fileext); res {
				c.Ctx.Output.Download(fileurl + "/" + attachment.FileName)
			} else {
				c.Data["json"] = map[string]interface{}{"info": "ERROR", "data": "权限不够！"}
				c.ServeJSON()
			}
		}
	} else {
		c.Data["json"] = map[string]interface{}{"info": "ERROR", "data": "未查到openID"}
		c.ServeJSON()
	}
}

// @Title dowload wx pdf
// @Description get wx pdf by id
// @Param id path string  true "The id of pdf"
// @Success 200 {object} models.GetAttachbyId
// @Failure 400 Invalid page supplied
// @Failure 404 pdf not found
// @router /getwxpdf/:id [get]
// 查阅防腐资料，不用权限判断
func (c *AttachController) GetWxPdf() {
	id := c.Ctx.Input.Param(":id")
	//pid转成64为
	idNum, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		logs.Error(err)
	}

	//根据附件id取得附件的prodid，路径
	attachment, err := models.GetAttachbyId(idNum)
	if err != nil {
		logs.Error(err)
	}
	fileext := path.Ext(attachment.FileName)
	matched, err := regexp.MatchString("\\.*[m|M][c|C][d|D]", fileext)
	if err != nil {
		logs.Error(err)
	}
	// beego.Info(matched)
	if matched {
		c.Data["json"] = map[string]interface{}{"info": "ERROR", "data": "不能下载mcd文件!"}
		c.ServeJSON()
		return
	}

	product, err := models.GetProd(attachment.ProductId)
	if err != nil {
		logs.Error(err)
	}

	//由proj id取得url
	fileurl, _, err := GetUrlPath(product.ProjectId)
	if err != nil {
		logs.Error(err)
		c.Data["json"] = "文件路径不存在！"
		c.ServeJSON()
	} else {
		c.Ctx.Output.Download(fileurl + "/" + attachment.FileName)
	}
}

// @Title dowload math pdf
// @Description get math pdf by link
// @Param id path string  true "The id of pdf"
// @Success 200 {object} models.GetAttachbyId
// @Failure 400 Invalid page supplied
// @Failure 404 pdf not found
// @router /mathpdf/:id [get]
// 下载mathcad pdf计算书
// 重复下载不扣费
func (c *AttachController) MathPdf() {
	// 加权限判断
	var err error
	var uintid uint
	_, _, uid, _, islogin := checkprodRole(c.Ctx)
	if !islogin {
		c.Data["json"] = map[string]interface{}{"info": "ERROR", "state": "ERROR", "data": "用户未登录！"}
		c.ServeJSON()
		return
	}

	// pdflink := c.GetString("pdflink")
	id := c.Ctx.Input.Param(":id")
	if id != "" {
		//id转成uint为
		idint, err := strconv.Atoi(id)
		if err != nil {
			logs.Error(err)
		}
		uintid = uint(idint)
	}
	//根据history uint查出路径
	history, err := models.GetMathHistoryByID(uintid)
	if err != nil {
		logs.Error(err)
	}

	// beego.Info(pdflink)
	fileext := path.Ext(history.PdfUrl)
	filename := filepath.Base(history.PdfUrl)
	// 查询uid，pdftitle记录，如果有，则直接下载
	_, err = models.GetUserPayMathPdf(uid, filename)
	if err != nil {
		logs.Error(err)
		// beego.Info(userpaymathpdf)

		// 查询余额
		// money, err := models.GetUserMoney(uid)
		// if err != nil {
		// 	logs.Error(err)
		// 	c.Data["json"] = map[string]interface{}{"info": "ERROR", "data": "查询余额错误！"}
		// 	c.ServeJSON()
		// 	return
		// }
		// if money.Amount <= 0 {
		// 	c.Data["json"] = map[string]interface{}{"info": "ERROR", "data": "用户余额不足，请联系管理员充值！"}
		// 	c.ServeJSON()
		// 	return
		// }

		matched, err := regexp.MatchString("\\.*[m|M][c|C][d|D]", fileext)
		if err != nil {
			logs.Error(err)
		}
		if matched {
			c.Data["json"] = map[string]interface{}{"info": "ERROR", "data": "不能下载mcd文件！"}
			c.ServeJSON()
			return
		}
		// 账户扣除2元，下载pdf扣除2元
		// err = models.AddUserPayMathPdf(filename, history.PdfUrl, uid, 2)
		// if err != nil {
		// 	logs.Error(err)
		// 	c.Data["json"] = map[string]interface{}{"info": "ERROR", "data": "扣费发生错误！"}
		// 	c.ServeJSON()
		// 	return
		// }
	}
	// pdflink = strings.Replace(pdflink, "/attachment", "attachment", -1)
	// c.Ctx.Output.Download(pdflink)
	c.Data["PdfLink"] = history.PdfUrl // 这里应该用id
	c.TplName = "web/viewer.html"
}

// @Title dowload wx pdf
// @Description get wx pdf by link
// @Param id path string true "The id of history pdf"
// @Success 200 {object} models.GetAttachbyId
// @Failure 400 Invalid page supplied
// @Failure 404 pdf not found
// @router /wxmathpdf/:id [get]
// 下载mathcad pdf计算书
func (c *AttachController) WxMathPdf() {
	// 加权限判断
	var user models.User
	var err error
	var uintid uint
	openID := c.GetSession("openID")
	if openID != nil {
		user, err = models.GetUserByOpenID(openID.(string))
		if err != nil {
			logs.Error(err)
			c.Data["json"] = map[string]interface{}{"info": "ERROR", "data": "获取用户名错误！"}
			c.ServeJSON()
			return
		}
	} else {
		c.Data["json"] = map[string]interface{}{"info": "ERROR", "data": "用户未登录！"}
		c.ServeJSON()
		return
	}

	// pdflink := c.GetString("pdflink")
	id := c.Ctx.Input.Param(":id")
	if id != "" {
		//id转成uint为
		idint, err := strconv.Atoi(id)
		if err != nil {
			logs.Error(err)
		}
		uintid = uint(idint)
	}
	//根据history uint查出路径
	history, err := models.GetMathHistoryByID(uintid)
	if err != nil {
		logs.Error(err)
	}

	// beego.Info(pdflink)
	fileext := path.Ext(history.PdfUrl)
	filename := filepath.Base(history.PdfUrl)

	// beego.Info(pdflink)
	// fileext := path.Ext(pdflink)
	// filename := filepath.Base(pdflink)
	// 查询uid，pdftitle记录，如果有，则直接下载
	_, err = models.GetUserPayMathPdf(user.Id, filename)
	if err != nil {
		logs.Error(err)

		// 查询余额
		// money, err := models.GetUserMoney(user.Id)
		// if err != nil {
		// 	logs.Error(err)
		// 	c.Data["json"] = map[string]interface{}{"info": "ERROR", "data": "查询余额错误！"}
		// 	c.ServeJSON()
		// 	return
		// }
		// if money.Amount <= 0 {
		// 	c.Data["json"] = map[string]interface{}{"info": "ERROR", "data": "用户余额不足，请联系管理员充值！"}
		// 	c.ServeJSON()
		// 	return
		// }

		matched, err := regexp.MatchString("\\.*[m|M][c|C][d|D]", fileext)
		if err != nil {
			logs.Error(err)
		}
		if matched {
			c.Data["json"] = map[string]interface{}{"info": "ERROR", "data": "不能下载mcd文件！"}
			c.ServeJSON()
			return
		}

		matched2, err := regexp.MatchString(".*[x|X][l|L][s|S]", fileext)
		if err != nil {
			logs.Error(err)
			c.Data["json"] = map[string]interface{}{"info": "ERROR", "data": "正则匹配错误！"}
			c.ServeJSON()
			return
		}
		// beego.Info(matched)
		if matched2 {
			c.Data["json"] = map[string]interface{}{"info": "ERROR", "data": "不能下载excel文件！"}
			c.ServeJSON()
			return
		}

		// 账户扣除2元，下载pdf扣除2元
		// err = models.AddUserPayMathPdf(filename, history.PdfUrl, user.Id, 2)
		// if err != nil {
		// 	logs.Error(err)
		// 	c.Data["json"] = map[string]interface{}{"info": "ERROR", "data": "扣费发生错误！"}
		// 	c.ServeJSON()
		// 	return
		// }
	}
	// 微信必须要下面这一步进行替换，否则无法下载，并且返回是404
	pdflink := strings.Replace(history.PdfUrl, "/attachment", "attachment", -1)
	c.Ctx.Output.Download(pdflink)
}

// @Title dowload wx math temp pdf
// @Description get wx math temp pdf by id
// @Param id path string true "The url of pdf"
// @Success 200 {object} models.GetAttachbyId
// @Failure 400 Invalid page supplied
// @Failure 404 pdf not found
// @router /wxmathtemppdf/:id [get]
// 下载mathcad模板 pdf阅览
func (c *AttachController) WxMathTempPdf() {
	// 加权限判断
	openID := c.GetSession("openID")
	if openID != nil {
		logs.Info(openID.(string))
		_, err := models.GetUserByOpenID(openID.(string))
		if err != nil {
			logs.Error(err)
			c.Data["json"] = map[string]interface{}{"info": "ERROR", "data": "获取用户名错误!"}
			c.ServeJSON()
			return
		}
	} else {
		hotqinsessionid := c.GetString("hotqinsessionid")
		// logs.Info(hotqinsessionid)
		//这里用security.go里的方法
		requestUrl := "http://127.0.0.1:8082/v1/wx/passlogin?hotqinsessionid=" + hotqinsessionid
		resp, err := http.Get(requestUrl)
		if err != nil {
			logs.Error(err)
			// return
		}
		defer resp.Body.Close()
		var data map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			logs.Error(err)
		}
		if _, ok := data["token"]; !ok {
			info := data["info"].(string)
			logs.Info(info)
			data := data["data"].(string)
			c.Data["json"] = map[string]interface{}{"info": info, "state": "ERROR", "data": data}
			c.ServeJSON()
			return
		} else {
			var openID, token string
			openID = data["openid"].(string)
			token = data["token"].(string)
			logs.Info(openID)
			logs.Info(token)
			// 对token进行验证
			userName, err := utils.CheckToken(token)
			if err != nil {
				logs.Error(err)
			}
			if userName == "" {
				c.Data["json"] = map[string]interface{}{"info": "ERROR", "state": "ERROR", "data": "用户未登录"}
				c.ServeJSON()
				return
			}
		}
	}

	id := c.Ctx.Input.Param(":id")
	var mathtempleid uint
	if id != "" {
		//id转成uint为
		idint, err := strconv.Atoi(id)
		if err != nil {
			logs.Error(err)
		}
		mathtempleid = uint(idint)
	}

	mathtemple, err := models.GetMathTemple(mathtempleid)
	if err != nil {
		logs.Error(err)
		c.Data["json"] = map[string]interface{}{"info": "ERROR", "state": "ERROR", "data": "查询数据库错误！"}
		c.ServeJSON()
		return
	}
	// beego.Info(mathtemple)
	// 去除文件名
	filepath := path.Dir(mathtemple.TempPath)
	// 文件名
	filename := mathtemple.TempTitle //path.Base(mathtemple.TempPath)
	// 文件后缀
	filesuffix := path.Ext(filename)
	filenameOnly := strings.TrimSuffix(filename, filesuffix) //只留下文件名，无后缀
	// fileext := path.Ext(pdflink)
	// matched, err := regexp.MatchString("\\.*[m|M][c|C][d|D]", fileext)
	// if err != nil {
	// 	logs.Error(err)
	// }
	// if matched {
	// 	c.Data["json"] = map[string]interface{}{"info": "ERROR", "data": "不能下载mcd文件!"}
	// 	c.ServeJSON()
	// 	return
	// }
	pdflink := filepath + "/" + filenameOnly + ".pdf"
	// logs.Info(pdflink)
	c.Ctx.Output.Download(pdflink)
}

// @Title dowload math pdf
// @Description get math pdf by link
// @Param id path string  true "The id of pdf"
// @Success 200 {object} models.GetAttachbyId
// @Failure 400 Invalid page supplied
// @Failure 404 pdf not found
// @router /excelpdf/:id [get]
// 下载excel pdf计算书
// 重复下载不扣费
func (c *AttachController) ExcelPdf() {
	// 加权限判断
	var err error
	var uintid uint
	_, _, uid, _, islogin := checkprodRole(c.Ctx)
	if !islogin {
		c.Data["json"] = map[string]interface{}{"info": "ERROR", "state": "ERROR", "data": "用户未登录！"}
		c.ServeJSON()
		return
	}

	// pdflink := c.GetString("pdflink")
	id := c.Ctx.Input.Param(":id")
	if id != "" {
		//id转成uint为
		idint, err := strconv.Atoi(id)
		if err != nil {
			logs.Error(err)
		}
		uintid = uint(idint)
	}
	//根据history uint查出路径
	history, err := models.GetHistoryExcel(uintid) //这里修改
	if err != nil {
		logs.Error(err)
	}

	// beego.Info(pdflink)
	fileext := path.Ext(history.PdfUrl)
	filename := filepath.Base(history.PdfUrl)
	// 查询uid，pdftitle记录，如果有，则直接下载
	_, err = models.GetUserPayExcelPdf(uid, filename)
	if err != nil {
		logs.Error(err)
		// beego.Info(userpaymathpdf)

		// 查询余额
		money, err := models.GetUserMoney(uid)
		if err != nil {
			logs.Error(err)
			c.Data["json"] = map[string]interface{}{"info": "ERROR", "data": "查询余额错误！"}
			c.ServeJSON()
			return
		}
		if money.Amount <= 0 {
			c.Data["json"] = map[string]interface{}{"info": "ERROR", "data": "用户余额不足，请联系管理员充值！"}
			c.ServeJSON()
			return
		}

		matched, err := regexp.MatchString("\\.*[m|M][c|C][d|D]", fileext)
		if err != nil {
			logs.Error(err)
		}
		if matched {
			c.Data["json"] = map[string]interface{}{"info": "ERROR", "data": "不能下载mcd文件！"}
			c.ServeJSON()
			return
		}
		// 账户扣除2元，下载pdf扣除2元
		err = models.AddUserPayExcelPdf(filename, history.PdfUrl, uid, 2)
		if err != nil {
			logs.Error(err)
			c.Data["json"] = map[string]interface{}{"info": "ERROR", "data": "扣费发生错误！"}
			c.ServeJSON()
			return
		}
	}
	// pdflink = strings.Replace(pdflink, "/attachment", "attachment", -1)
	// c.Ctx.Output.Download(pdflink)
	c.Data["PdfLink"] = history.PdfUrl // 这里应该用id
	c.TplName = "web/viewer.html"
}

// @Title dowload wx pdf
// @Description get wx pdf by link
// @Param id query string  true "The url of pdf"
// @Success 200 {object} models.GetAttachbyId
// @Failure 400 Invalid page supplied
// @Failure 404 pdf not found
// @router /wxexcelpdf/:id [get]
// 小程序端下载mexcel pdf计算书
func (c *AttachController) WxExcelPdf() {
	// 加权限判断
	var user models.User
	var err error
	openID := c.GetSession("openID")
	if openID != nil {
		user, err = models.GetUserByOpenID(openID.(string))
		if err != nil {
			logs.Error(err)
			c.Data["json"] = map[string]interface{}{"info": "ERROR", "data": "获取用户名错误！"}
			c.ServeJSON()
			return
		}
	} else {
		c.Data["json"] = map[string]interface{}{"info": "ERROR", "data": "用户未登录！"}
		c.ServeJSON()
		return
	}

	// pdflink := c.GetString("pdflink")
	var uintid uint
	id := c.Ctx.Input.Param(":id")
	if id != "" {
		//id转成uint为
		idint, err := strconv.Atoi(id)
		if err != nil {
			logs.Error(err)
		}
		uintid = uint(idint)
	}
	//根据history uint查出路径
	history, err := models.GetHistoryExcel(uintid)
	if err != nil {
		logs.Error(err)
	}

	// beego.Info(pdflink)
	fileext := path.Ext(history.PdfUrl)
	filename := filepath.Base(history.PdfUrl)
	// beego.Info(pdflink)
	// fileext := path.Ext(pdflink)
	// filename := filepath.Base(pdflink)
	// 查询uid，pdftitle记录，如果有，则直接下载
	_, err = models.GetUserPayExcelPdf(user.Id, filename)
	if err != nil {
		logs.Error(err)

		// 查询余额
		money, err := models.GetUserMoney(user.Id)
		if err != nil {
			logs.Error(err)
			c.Data["json"] = map[string]interface{}{"info": "ERROR", "data": "查询余额错误！"}
			c.ServeJSON()
			return
		}
		if money.Amount <= 0 {
			c.Data["json"] = map[string]interface{}{"info": "ERROR", "data": "用户余额不足，请联系管理员充值！"}
			c.ServeJSON()
			return
		}

		matched, err := regexp.MatchString("\\.*[m|M][c|C][d|D]", fileext)
		if err != nil {
			logs.Error(err)
		}
		if matched {
			c.Data["json"] = map[string]interface{}{"info": "ERROR", "data": "不能下载mcd文件！"}
			c.ServeJSON()
			return
		}

		matched2, err := regexp.MatchString(".*[x|X][l|L][s|S]", fileext)
		if err != nil {
			logs.Error(err)
			c.Data["json"] = map[string]interface{}{"info": "ERROR", "data": "正则匹配错误！"}
			c.ServeJSON()
			return
		}
		// beego.Info(matched)
		if matched2 {
			c.Data["json"] = map[string]interface{}{"info": "ERROR", "data": "不能下载excel文件！"}
			c.ServeJSON()
			return
		}

		// 账户扣除2元，下载pdf扣除2元
		err = models.AddUserPayExcelPdf(filename, history.PdfUrl, user.Id, 2)
		if err != nil {
			logs.Error(err)
			c.Data["json"] = map[string]interface{}{"info": "ERROR", "data": "扣费发生错误！"}
			c.ServeJSON()
			return
		}
	}
	// logs.Info(history.PdfUrl)
	pdflink := strings.Replace(history.PdfUrl, "/attachment", "attachment", -1)
	c.Ctx.Output.Download(pdflink)
}

// @Title dowload wx math temp pdf
// @Description get wx math temp pdf by id
// @Param id path string true "The url of pdf"
// @Success 200 {object} models.GetAttachbyId
// @Failure 400 Invalid page supplied
// @Failure 404 pdf not found
// @router /getwxexceltemppdf/:id [get]
// 小程序下载excel模板 pdf阅览
func (c *AttachController) WxExcelTempPdf() {
	// 加权限判断
	openID := c.GetSession("openID")
	if openID != nil {
		// beego.Info(openID.(string))
		_, err := models.GetUserByOpenID(openID.(string))
		if err != nil {
			logs.Error(err)
			c.Data["json"] = map[string]interface{}{"info": "ERROR", "data": "获取用户名错误!"}
			c.ServeJSON()
			return
		}
	} else {
		c.Data["json"] = map[string]interface{}{"info": "ERROR", "data": "用户未登录"}
		c.ServeJSON()
		return
	}

	id := c.Ctx.Input.Param(":id")
	var mathtempleid uint
	if id != "" {
		//id转成uint为
		idint, err := strconv.Atoi(id)
		if err != nil {
			logs.Error(err)
		}
		mathtempleid = uint(idint)
	}

	mathtemple, err := models.GetExcelTemple(mathtempleid)
	if err != nil {
		logs.Error(err)
	}
	// beego.Info(mathtemple)
	// 去除文件名
	filepath := path.Dir(mathtemple.Path)
	// 文件名
	filename := mathtemple.Title //path.Base(mathtemple.TempPath)
	// 文件后缀
	filesuffix := path.Ext(filename)
	filenameOnly := strings.TrimSuffix(filename, filesuffix) //只留下文件名，无后缀
	// fileext := path.Ext(pdflink)
	// matched, err := regexp.MatchString("\\.*[m|M][c|C][d|D]", fileext)
	// if err != nil {
	// 	logs.Error(err)
	// }
	// if matched {
	// 	c.Data["json"] = map[string]interface{}{"info": "ERROR", "data": "不能下载mcd文件!"}
	// 	c.ServeJSON()
	// 	return
	// }
	pdflink := filepath + "/" + filenameOnly + ".pdf"
	// beego.Info(pdflink)
	c.Ctx.Output.Download(pdflink)
}

// @Title dowload ansys data
// @Description get ansys data by link
// @Param id path string  true "The id of data"
// @Success 200 {object} models.GetAttachbyId
// @Failure 400 Invalid page supplied
// @Failure 404 data not found
// @router /ansysdata/:id [get]
// ansys data数据
// 重复下载不扣费
func (c *AttachController) AnsysData() {
	// 加权限判断
	var err error
	var uintid uint
	_, _, uid, isadmin, islogin := checkprodRole(c.Ctx)
	if !islogin {
		c.Data["json"] = map[string]interface{}{"info": "ERROR", "state": "ERROR", "data": "用户未登录！"}
		c.ServeJSON()
		return
	}

	id := c.Ctx.Input.Param(":id")
	// logs.Info(id)
	if id != "" {
		//id转成uint为
		idint, err := strconv.Atoi(id)
		if err != nil {
			logs.Error(err)
		}
		uintid = uint(idint)
	}

	//根据history uint查出路径
	history, err := models.GetHistoryAnsys(uintid) //这里修改
	if err != nil {
		logs.Error(err)
	}

	var isme bool
	if history.UserID == uid {
		isme = true
	}
	if isme || isadmin {
		// logs.Info(history)
		// fileext := path.Ext(history.PdfUrl)
		filename := filepath.Base(history.PdfUrl)
		// logs.Info(filename)
		// beego.Info(mathtemple)
		// 去除文件名
		filepath := path.Dir(history.PdfUrl)
		// 文件名
		// filename := mathtemple.Title //path.Base(mathtemple.TempPath)
		// 文件后缀
		filesuffix := path.Ext(history.PdfUrl)
		filenameOnly := strings.TrimSuffix(filename, filesuffix) //只留下文件名，无后缀
		// logs.Info(filenameOnly)
		pdflink := filepath + "/" + filenameOnly + ".dat"
		pdflink = strings.Replace(pdflink, "/attachment", "attachment", -1)
		// logs.Info(pdflink)
		c.Ctx.Output.Download(pdflink) //这句和下面都行，调试的时候不知道为什么这句一直出错！！
		// http.ServeFile(c.Ctx.ResponseWriter, c.Ctx.Request, pdflink)
	} else {
		c.Data["json"] = map[string]interface{}{"info": "ERROR", "state": "ERROR", "data": "只能本人查阅！"}
		c.ServeJSON()
	}
}

//编码转换
// l3, err3 := url.Parse(c.Ctx.Request.RequestURI[1:])
// 	if err3 != nil {
// 		logs.Error(err3)
// 	}
// 	beego.Info(l3.Query())
// 	beego.Info(l3.Query().Encode())
// 	beego.Info(url.QueryEscape("初步设计"))
// 	dec := mahonia.NewDecoder("utf8") //定义转换乱码
// 	//保留	fmt.Println(dec.ConvertString(j))   //转成utf-8
// 	beego.Info(dec.ConvertString(url.QueryEscape("初步设计")))
// 	// rd := dec.NewReader(resp.Body)
// 	// doc, err := goquery.NewDocumentFromReader(rd)
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	bn := []byte("\xbf\xc6\xd1\xa7\xC3\xF1\xD6\xF7\xCF\xDC\xD5\xFE")
// 	beego.Info(bn)
// 	sn := string(bn)
// 	beego.Info(sn)
// 	cbn, err1, _, _ := ConvertGB2312(bn)
// 	if err1 != nil {
// 		logs.Error(err1)
// 	}
// 	beego.Info(cbn)
// 	csn, err2, _, _ := ConvertGB2312String(c.Ctx.Request.RequestURI[1:])
// 	if err2 != nil {
// 		logs.Error(err2)
// 	}
// 	beego.Info(csn)
// 	//根据路由path.Dir——再转成数组strings.Split——查出项目id——加上名称——查出下级id
// 	beego.Info(path.Dir(filePath))
// 	beego.Info(c.Ctx.Request.RequestURI[1:])

// 运算符优先级

// 由上至下代表优先级由高到低

// 7	^ !
// 6	* / % << >> & &^
// 5	+ - | ^
// 4	== != < <= >= >
// 3	<-
// 2	&&
// 1	||
