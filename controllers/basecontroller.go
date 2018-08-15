// Copyright 2018. box.la authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controllers

import (
	"github.com/astaxie/beego"
)

type baseController struct {
	beego.Controller
}

func (this *baseController) Prepare() {
	//签名验证等
	//this.Data["json"] = map[string]interface{}{"msg": "verify not pass"}
	//this.ServeJSON()
	//this.StopRun()
}
