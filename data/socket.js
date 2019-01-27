///handle enter pressed
var input = document.getElementById("command_input");
input.addEventListener("keyup", function(event) {
  event.preventDefault();
  if (event.keyCode === 13) {
    document.getElementById("submit_command").click();
  }
});


var login_token = "dedtoken";
var had_connected = false;
var address = "46.101.66.98:8080";
var output = document.getElementById("output");
var socket = new WebSocket("ws://"+address+"/ws");

socket.onopen = function () {
	print_entry("Successfully connected to " + address);
	had_connected = true;
	socket.send("new_conn");
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
	var token_found = false;
	//waiting for new token to arrive
	if(login_token == "dedtoken")
	{
		var split_tokens = e.data.split(":");
		if(split_tokens.length == 2 && split_tokens[0] == "userid")
		{
			login_token = split_tokens[1];
			token_found = true;
			print_entry("Logged in successfully!");	
		}
	}
	if(token_found == false)
	{
		print_entry(e.data)
	}
};

function send() {
	if(input.value == "clear")
	{
		output.innerHTML = "";
	}
	else if(input.value != "")
	{
		socket.send(login_token + " " + input.value);	
	}
	input.value = "";
}

function print_entry(str){
	output.innerHTML += "<div class='entry'>"+str+"</div>"
	window.scrollTo(0,document.body.scrollHeight);
}