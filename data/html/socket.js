///handle enter pressed
var input = document.getElementById("command_input");
input.addEventListener("keyup", function(event) {
  event.preventDefault();
  if (event.keyCode === 13) {
    document.getElementById("submit_command").click();
  }
});

var login_token = "";
var had_connected = false;
var address = window.location.hostname+":8080";
var output = document.getElementById("output");
var socket = new WebSocket("ws://"+address+"/ws");

socket.onopen = function () {
	print_entry("Successfully connected to " + address);
	had_connected = true;
	send_with_token("new_conn");
};

socket.onclose = function () {
	if(had_connected)
	{
		print_entry("Connection to " + address + " was lost.");	
	}
	else
	{
		print_entry("Failed to connect to " + address);
	}
	
};

socket.onmessage = function (e) {
	var recieved_obj = JSON.parse(e.data);
	if(recieved_obj["token"] != "")
	{
		login_token = recieved_obj["token"]
	}

	if(recieved_obj["message"] != "")
	{
		print_entry(recieved_obj["message"]);
	}
};

function send() {
	if(input.value == "clear")
	{
		//client side clear
		output.innerHTML = "";
	}
	else if(input.value != "")
	{
		send_with_token(input.value);
	}
	input.value = "";
}

function send_with_token(str)
{
	var to_send = {
		"token" : login_token,
		"message" : str
	};
	socket.send(JSON.stringify(to_send));
}

function print_entry(str){
	output.innerHTML += "<div class='entry'>"+str+"</div>"
	window.scrollTo(0,document.body.scrollHeight);
}