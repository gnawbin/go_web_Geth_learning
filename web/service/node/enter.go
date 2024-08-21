package node

import "k8s-web/convert"

// @Author: morris
type Group struct {
	NodeService
}

var nodeConvert = convert.ConvertGroupApp.NodeConvert
