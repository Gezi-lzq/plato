package config

import "github.com/spf13/viper"

func GetStateCmdChannelNum() int {
	return viper.GetInt("state.cmd_channel_num")
}
func GetStateServiceAddr() string {
	return viper.GetString("state.servide_addr")
}
func GetStateServiceName() string {
	return viper.GetString("state.service_name")
}
func GetStateServerPort() int {
	return viper.GetInt("state.server_port")
}
func GetStateRPCWeight() int {
	return viper.GetInt("state.weight")
}
func GetStateRPCProtocol() string {
	return viper.GetString("state.rpc_protocol")
}
func GetStateRPCSockAdd() string {
	return viper.GetString("state.rpc_sockAddr")
}
