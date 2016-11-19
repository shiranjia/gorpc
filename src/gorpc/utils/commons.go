package utils

import "log"

func CheckErr(e error) error {
	if e!= nil{
		log.Println("Err:",e)
	}
	return e
}

func HandlerErr(e error , handler func(interface{}))  {
	if e!= nil{
		log.Println("Err:",e)
		if handler != nil{
			handler(e)
		}
	}
}

func Try(do func(),handler func(interface{}))  {
	defer func(){
		if err := recover();err != nil{
			handler(err)
		}
	}()
	do()
}
