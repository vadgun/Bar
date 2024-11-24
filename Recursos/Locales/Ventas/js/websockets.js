$(document).ready(function() {
    var socket = new WebSocket("ws://192.168.1.67:8080/ws");

    //Al abrir la vista venta.html
    socket.onopen = function() {
        console.log("Connected to WebSocket server");
    };

    //Al propagarse un mensaje
    socket.onmessage = function(event) {
        // var message = JSON.parse(event.data);
        // $('#messages').append('<p><strong>' + message.username + ':</strong> ' + message.message + '</p>');

        // const data = JSON.parse(event.data);
        //     console.log(data);
        //     document.getElementById('maincontainer').className = 'newClass';
        VerificaMesas(getCurrentDate());
        MifuncionModal($("#tablaactiva").val());
        if ($("#mesaenmodal").val() != undefined) {
            ProductosYaVendidos($("#mesaenmodal").val());
          }
    };

    //Al cerrar la ventana del navegador venta, o cambiar de pestana
    socket.onclose = function() {
        console.log("Disconnected from WebSocket server");
    };

    //Al hacer click al boton de send
    // $('#send').click(function() {
    //     var username = $('#username').val();
    //     var message = $('#message').val();
    //     socket.send(JSON.stringify({ username: username, message: message }));
    //     $('#message').val('');
    // });

});

//Construye la fecha en el siguiente formato 2024-07-23
const getCurrentDate = () => {
    const date = new Date();
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    return `${year}-${month}-${day}`;
};

$(document).ready(function() {
    VerificaMesas(getCurrentDate());
})