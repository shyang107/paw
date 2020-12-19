package main

func exlog() {
	log.Infoln("飛雪無情的博客:", "http://www.flysnow.org")
	log.Warnf("飛雪無情的微信公眾號：%s\n", "flysnow_org")
	log.Errorln("歡迎關注留言")

	lg.Infoln("飛雪無情的博客:", "http://www.flysnow.org")
	lg.Warnln("飛雪無情的博客:", "http://www.flysnow.org")
	lg.Debugln("飛雪無情的博客:", "http://www.flysnow.org")
	lg.Errorln("飛雪無情的博客:", "http://www.flysnow.org")
	// lg.Traceln("飛雪無情的博客:", "http://www.flysnow.org")
	// lg.Fatalln("飛雪無情的博客:", "http://www.flysnow.org")

}
