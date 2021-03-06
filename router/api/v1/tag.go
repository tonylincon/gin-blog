package v1

import (
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"

	"github.com/linehk/gin-blog/errno"
	"github.com/linehk/gin-blog/router/api"
	"github.com/linehk/gin-blog/vm"
)

func GetTags(c *gin.Context) {
	// 返回 URL 参数的值
	name := c.Query("name")
	state := -1
	if s := c.Query("state"); s != "" {
		state = com.StrTo(s).MustInt()
	}

	// 构造结构体
	vmTag := vm.Tag{Name: name, State: state, PageNum: api.PageNum(c), PageSize: api.PageSize}

	tags, err := vmTag.GetAll()
	if err != nil {
		api.Response(c, http.StatusInternalServerError, errno.ERROR_GET_TAGS_FAIL, nil)
		return
	}

	// 计数
	count, err := vmTag.Count()
	if err != nil {
		api.Response(c, http.StatusInternalServerError, errno.ERROR_COUNT_TAG_FAIL, nil)
		return
	}

	// 填充数据
	data := map[string]interface{}{"lists": tags, "count": count}

	api.Response(c, http.StatusOK, errno.SUCCESS, data)
}

func GetTag(c *gin.Context) {
	// 获取 id
	id := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(id, 1, "id")

	// 表单验证错误
	if valid.HasErrors() {
		api.LogErrors(valid.Errors)
		api.Response(c, http.StatusBadRequest, errno.INVALID_PARAMS, nil)
	}

	vmTag := vm.Tag{ID: id}
	exist, err := vmTag.HasID()
	if err != nil {
		api.Response(c, http.StatusInternalServerError, errno.ERROR_GET_TAGS_FAIL, nil)
		return
	}
	if !exist {
		api.Response(c, http.StatusOK, errno.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	tag, err := vmTag.Get()
	if err != nil {
		api.Response(c, http.StatusInternalServerError, errno.ERROR_GET_ARTICLE_FAIL, nil)
		return
	}
	api.Response(c, http.StatusOK, errno.SUCCESS, tag)
}

func HasTagByName(c *gin.Context) {
}

func HasTagByID(c *gin.Context) {
}

func GetTagsCount(c *gin.Context) {
}

type AddTagForm struct {
	Name      string `form:"name" valid:"Required;MaxSize(100)"`
	CreatedBy string `form:"created_by" valid:"Required;MaxSize(100)"`
	State     int    `form:"state" valid:"Range(0,1)"`
}

func AddTag(c *gin.Context) {
	var form AddTagForm
	httpCode, errCode := api.BindAndValid(c, &form)
	if errCode != errno.SUCCESS {
		api.Response(c, httpCode, errCode, nil)
		return
	}

	vmTag := vm.Tag{Name: form.Name, CreatedBy: form.CreatedBy, State: form.State}

	exist, err := vmTag.HasName()
	if err != nil {
		api.Response(c, http.StatusInternalServerError, errno.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if exist {
		api.Response(c, http.StatusOK, errno.ERROR_EXIST_TAG, nil)
		return
	}

	err = vmTag.Add()
	if err != nil {
		api.Response(c, http.StatusInternalServerError, errno.ERROR_ADD_TAG_FAIL, nil)
		return
	}

	api.Response(c, http.StatusOK, errno.SUCCESS, nil)
}

type EditTagForm struct {
	ID         int    `form:"id" valid:"Required;Min(1)"`
	Name       string `form:"name" valid:"Required;MaxSize(100)"`
	ModifiedBy string `form:"modified_by" valid:"Required;MaxSize(100)"`
	State      int    `form:"state" valid:"Range(0,1)"`
}

func EditTag(c *gin.Context) {
	form := EditTagForm{ID: com.StrTo(c.Param("id")).MustInt()}
	httpCode, errCode := api.BindAndValid(c, &form)
	if errCode != errno.SUCCESS {
		api.Response(c, httpCode, errCode, nil)
		return
	}

	vmTag := vm.Tag{ID: form.ID, Name: form.Name, ModifiedBy: form.ModifiedBy, State: form.State}

	exist, err := vmTag.HasID()
	if err != nil {
		api.Response(c, http.StatusInternalServerError, errno.ERROR_EXIST_TAG_FAIL, nil)
		return
	}

	if !exist {
		api.Response(c, http.StatusOK, errno.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	err = vmTag.Edit()
	if err != nil {
		api.Response(c, http.StatusInternalServerError, errno.ERROR_EDIT_TAG_FAIL, nil)
		return
	}
	api.Response(c, http.StatusOK, errno.SUCCESS, nil)
}

func DeleteTag(c *gin.Context) {
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt()
	valid.Min(id, 1, "id").Message("ID必须大于0")

	if valid.HasErrors() {
		api.LogErrors(valid.Errors)
		api.Response(c, http.StatusBadRequest, errno.INVALID_PARAMS, nil)
	}

	vmTag := vm.Tag{ID: id}
	exist, err := vmTag.HasID()
	if err != nil {
		api.Response(c, http.StatusInternalServerError, errno.ERROR_EXIST_TAG_FAIL, nil)
		return
	}

	if !exist {
		api.Response(c, http.StatusOK, errno.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	if err := vmTag.Delete(); err != nil {
		api.Response(c, http.StatusInternalServerError, errno.ERROR_DELETE_TAG_FAIL, nil)
		return
	}

	api.Response(c, http.StatusOK, errno.SUCCESS, nil)
}

func DeleteTags(c *gin.Context) {
}
