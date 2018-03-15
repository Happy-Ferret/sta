package display

// Motd writes a hello line to the display.
func (d *Display) Motd() {
	d.WriteLine(`
 #####   #######     #     
#     #     #       # #    
#           #      #   #   
 #####      #     #     #  
      #     #     #######  
#     #     #     #     #  
 #####      #     #     #
`)
}
