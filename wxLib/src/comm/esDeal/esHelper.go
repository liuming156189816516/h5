package esDeal

import (
	"comm/comm"
	"context"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	elastic "gopkg.in/olivere/elastic.v7"
)

type EsConfig struct {
	Sniff     bool
	Url       string
	User      string
	Pwd       string
	IndexName []string
}

type DataInfo struct {
	Phone           string `json:"phone,omitempty"`              // 手机号
	PhoneMd5        string `json:"phone_md5,omitempty"`          // 手机号MD5
	Wxid            string `json:"wxid,omitempty"`               // 微信ID
	Alias           string `json:"alias,omitempty"`              // 微信号
	RealName        string `json:"real_name,omitempty"`          // 微信实名昵称
	RealState       int64  `json:"real_state,omitempty"`         // 微信实名状态 -1 不存在 1 存在
	Sex             int64  `json:"sex,omitempty"`                // 微信性别
	NickName        string `json:"nick_name,omitempty"`          // 微信昵称
	Signature       string `json:"signature,omitempty"`          // 微信签名
	Country         string `json:"country,omitempty"`            // 国家
	Province        string `json:"province,omitempty"`           // 微信省份
	City            string `json:"city,omitempty"`               // 微信城市
	SmallHeadImgUrl string `json:"small_head_img_url,omitempty"` // 头像地址
	ZfbRealName     string `json:"zfb_real_name,omitempty"`      // 支付宝实名昵称
	ZfbRealState    int64  `json:"zfb_real_state,omitempty"`     // 支付宝实名状态 -1 不存在 1 存在
	ZfbSex          int64  `json:"zfb_sex,omitempty"`            // 支付宝别
	ZfbNickName     string `json:"zfb_nick_name,omitempty"`      // 支付宝昵称
	ZfbProvince     string `json:"zfb_province,omitempty"`       // 支付宝省份
	ZfbCity         string `json:"zfb_city,omitempty"`           // 支付宝城市
	Constellation   string `json:"constellation,omitempty"`      // 星座
	UpTime          int64  `json:"up_time,omitempty"`            // 最后更新时间
}

type QqDataInfo struct {
	Phone    string `json:"phone,omitempty"`     // 手机号
	QqNumber string `json:"qq_number,omitempty"` // QQ号
}

//创建索引库时以下可以不指定，走系统默认类型；防止分词查询问题，这里直接创建，指定映射，做初始化使用
const wx_mapping = `{ 
"mappings":{
    "properties":{
        "phone":{"type":"keyword"},
        "phone_md5":{"type":"keyword"},
        "wxid":{"type":"keyword"},
        "alias":{"type":"keyword"},
        "real_name":{"type":"keyword"},
        "real_state":{"type":"integer"},
        "sex":{"type":"integer"},
        "nick_name":{"type":"keyword"},
        "signature":{"type":"text"},
        "country":{"type":"keyword"},
        "province":{"type":"keyword"},
        "city":{"type":"keyword"},
        "small_head_img_url":{"type":"text"},
        "zfb_real_name":{"type":"keyword"},
        "zfb_real_state":{"type":"integer"},
        "zfb_sex":{"type":"integer"},
        "zfb_nick_name":{"type":"keyword"},
        "zfb_province":{"type":"keyword"},
        "zfb_city":{"type":"keyword"},
        "constellation":{"type":"keyword"},
        "up_time":{"type":"long"}
    }
 }
}`

const qq_mapping = `{ 
"mappings":{
    "properties":{
        "phone":{"type":"keyword"},
        "qq_number":{"type":"keyword"},
    }
 }
}`

var (
	ctx     = context.Background()
	_client = &elastic.Client{}
)

// 获取客户端
func InitElastic(_esConfig *EsConfig) error {
	var err error
	options := []elastic.ClientOptionFunc{}
	options = append(options, elastic.SetSniff(_esConfig.Sniff))
	options = append(options, elastic.SetURL(_esConfig.Url, _esConfig.Url))
	if _esConfig.User != "" && _esConfig.Pwd != "" {
		options = append(options, elastic.SetBasicAuth(_esConfig.User, _esConfig.Pwd))
	}
	_client, err = elastic.NewClient(options...)
	if err != nil {
		logs.Error(fmt.Sprintf("get client config:%+s, err:%+v", _esConfig, err))
		return err
	}
	// 微信实名数据初始化 只会初始一次
	exists, err := _client.IndexExists(comm.GetFirstIndexName()).Do(ctx)
	if err != nil {
		return err
	}
	if !exists {
		logs.Debug(comm.GetFirstIndexName() + "--索引库不存在，进入初始化。。。。。。。")
		createIndex, errIndex := _client.CreateIndex(comm.GetFirstIndexName()).Body(wx_mapping).Do(ctx)
		if errIndex != nil {
			logs.Error(fmt.Sprintf("创建索引库失败---->index: %s, err:%s", comm.GetFirstIndexName(), errIndex.Error()))
			return errIndex
		}
		if !createIndex.Acknowledged {
			logs.Error(fmt.Sprintf("创建索引库失败,更新集群状态超时---->index: %s", comm.GetFirstIndexName()))
			return errors.New("创建索引库失败,更新集群状态超时")
		} else {
			logs.Debug(fmt.Sprintf("创建索引库成功---->index: %s", comm.GetFirstIndexName()))
		}
	} else {
		logs.Debug(comm.GetFirstIndexName() + "--索引库已存在。。。。。。。")
	}
	//检测
	info, code, err := _client.Ping(_esConfig.Url).Do(ctx)
	if err != nil {
		return err
	}
	logs.Debug(fmt.Sprintf("config:%+v, info:%+v, code:%d", _esConfig, info, code))

	version, err := _client.ElasticsearchVersion(_esConfig.Url)
	if err != nil {
		return err
	}
	logs.Trace("Elasticsearch version %s\n", version)
	return err
}

/**
插入
*/
func EsInsertInterface(dbName, TbName string, data interface{}) error {
	if _client == nil {
		return errors.New("Elasticsearch client is nil")
	}
	err := CreateIndex(dbName, "")
	if err != nil {
		return err
	}
	rsp, err := _client.Index().Index(dbName).Type(TbName).BodyJson(data).Do(ctx)
	if err != nil {
		return err
	}
	logs.Debug(fmt.Sprintf("rsp:%+v", rsp))
	return nil
}

func CreateIndex(index, mapping string) error {
	// 判断索引是否存在
	exists, err := _client.IndexExists(index).Do(ctx)
	if err != nil {
		logs.Error(fmt.Sprintf("<CreateIndex> some error occurred when check exists, index: %s, err:%s", index, err.Error()))
		return err
	}
	if exists {
		logs.Debug(fmt.Sprintf("<CreateIndex> index:{%s} is already exists", index))
		return nil
	}
	createIndex, err := _client.CreateIndex(index).Body(mapping).Do(ctx)
	if err != nil {
		logs.Error(fmt.Sprintf("<CreateIndex> some error occurred when create. index: %s, err:%s", index, err.Error()))
		return err
	}
	if !createIndex.Acknowledged {
		return errors.New(fmt.Sprintf("<CreateIndex> Not acknowledged, index: %s", index))
	}
	return nil
}

/**
删除索引
*/
func EsDeleteIndex(dbName string) error {
	rsp, err := _client.DeleteIndex(dbName).Do(ctx)
	if err != nil {
		return err
	}
	logs.Debug(fmt.Sprintf("rsp:%+v", rsp))
	return nil
}

/**
修改
*/
func EsUpdateById(dbName, TbName string, id string, data map[string]interface{}) error {
	if _client == nil {
		return errors.New("Elasticsearch client is nil")
	}
	rsp, err := _client.Update().Index(dbName).Type(TbName).Id(id).Doc(data).Do(ctx)
	if err != nil {
		return err
	}
	logs.Debug(fmt.Sprintf("rsp:%+v", rsp))
	return nil
}

/**
搜索
*/
func EsSearchPage(dbName, TbName string, query elastic.Query, page, pageSize int, sort, order string) (*elastic.SearchResult, error) {
	if _client == nil {
		return nil, errors.New("Elasticsearch client is nil")
	}

	result, err := _client.Search().Index(dbName).Type(TbName).Query(query).Pretty(true).Size(pageSize).From((page-1)*pageSize).Sort(sort, order == "ASC" || order == "asc").Do(ctx)
	if err != nil {
		return nil, err
	}

	return result, nil
}

/**
搜索
*/
func EsSearchALL(dbName, TbName string, query elastic.Query, size int) (*elastic.SearchResult, error) {
	if _client == nil {
		return nil, errors.New("Elasticsearch client is nil")
	}
	result, err := _client.Search().Index(dbName).Type(TbName).Query(query).Size(size).Do(ctx)
	if err != nil {
		return nil, err
	}
	return result, nil
}

/**
搜索
*/
func EsSearchById(dbName, TbName string, id string) (*elastic.GetResult, error) {
	if _client == nil {
		return nil, errors.New("Elasticsearch client is nil")
	}

	result, err := _client.Get().Index(dbName).Type(TbName).Id(id).Do(ctx)
	if err != nil {
		return nil, err
	}

	return result, nil
}

/**
删除
*/
func EsDeleteById(dbName, TbName string, id string) error {
	if _client == nil {
		return errors.New("Elasticsearch client is nil")
	}
	res, err := _client.Delete().Index(dbName).Type(TbName).Id(id).Do(ctx)
	if err != nil {
		return err
	}
	logs.Debug(fmt.Sprintf("rsp:%+v", res))
	return nil
}

/*------------7.x------------*/
/*插入or更新*/
func Upsert(indexName, id string, data interface{}) error {
	if _client == nil {
		return errors.New("Elasticsearch client is nil")
	}
	exists, err := _client.Exists().Index(indexName).Id(id).Do(context.TODO())
	if err != nil {
		logs.Error(fmt.Sprintf("Exists data Failed---->%s", err.Error()))
		return errors.New("Exists data Failed")
	}
	//空值不入库
	if exists { //非覆盖式更新
		_, err := _client.Update().Index(indexName).Id(id).Doc(data).Do(ctx)
		if err != nil {
			logs.Error(fmt.Sprintf("Update failed. index:%s,  id:%s err:%s", indexName, id, err.Error()))
			return err
		}
		//logs.Debug(fmt.Sprintf("Update successfully index:%s, id:%s, version:%d", indexName, res.Id, res.Version))
	} else { //插入
		_, err := _client.Index().Index(indexName).BodyJson(data).Id(id).Do(ctx)
		if err != nil {
			logs.Error(fmt.Sprintf("Insert data failed---->%s", err.Error()))
			return err
		}
		//logs.Debug(fmt.Sprintf("Insert data successfully---->:%+v", res))
	}
	return nil
}

/*查询*/
func Get(index, id string) (*elastic.GetResult, error) {
	if _client == nil {
		return nil, errors.New("Elasticsearch client is nil")
	}
	res, err := _client.Get().Index(index).Id(id).Do(ctx)
	if err != nil {
		//logs.Error(fmt.Sprintf("Get data Failed---->%s", err.Error()))
		return nil, err
	}
	return res, nil
}
