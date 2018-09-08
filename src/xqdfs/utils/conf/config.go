package conf

type Log struct {
	Level string		`json:"level"`
}

type Configure struct {
	Param string		`json:"param"`
}

type Server struct {
	Id int			`json:"id"`
	Desc string		`json:"desc"`
	Host string		`json:"host"`
	Port int		`json:"port"`
}
