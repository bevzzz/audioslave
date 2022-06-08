package main

//
//func main() {
//	keyEvents, err := keyboard.GetKeys(10)
//	if err != nil {
//		panic(err)
//	}
//	defer keyboard.Close()
//
//	volume := NewVolume()
//	for {
//		event := <-keyEvents
//		if event.Err != nil {
//			panic(err)
//		}
//		volume.Adjust()
//
//		state := "Off"
//		if volume.On() {
//			state = "On"
//		}
//		fmt.Println("Sound", state)
//	}
//}
