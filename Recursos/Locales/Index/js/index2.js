$(document).ready(function() {
    // alert("hola javascript");
    // $("#buscador").click(function(){
    //     alert("Hola jquery");
    // });

    // confirm dialog
    // alertify.confirm("¿ Cerrar sesión ?", function (e) {
    //     if (e) {
    //         alertify.success("Sing out");
    //     } else {
    //         alertify.error("Error notification");
    //     }
    // });
    $("#buscador").keyup(function() {
        buscar = $("#buscador").val();
        bs = buscar.toUpperCase();
        GetCoworker(bs);
    });


    // como prevenir la ida del formulario de alta de colaborador e implementar las mayusculas

})

function Descargar(name) {

    $.ajax({
        url: 'Recursos/Archivos/' + name + '.pdf',
        method: 'GET',
        xhrFields: {
            responseType: 'blob'
        },
        success: function(data) {
            var a = document.createElement('a');
            var url = window.URL.createObjectURL(data);
            a.href = url;
            a.download = name + '.pdf';
            a.click();
            window.URL.revokeObjectURL(url);
        }
    });
}


function Mifuncion(data) {
    console.log(data);
    $.ajax({
        url: '/getreport',
        data: { data: data },
        type: 'POST',
        dataType: 'html',
        success: function(result) {
            console.log("Operacion Realizada con Exito");
            $("#maincontainer").html(result);
        },
        error: function(xhr, status) {
            console.log("Error en la consulta")
        },
        complete: function(xhr, status) {
            console.log("Proceso Terminado")
        }
    });
}

function Antiguedades() {
    inicio = $("#inicio").val();
    fin = $("#fin").val();
    base = $("#base").val();

    if (base == "x") {
        alert("Selecciona base de datos")
        return false
    }

    $.ajax({
        url: '/antiguedadpagos',
        data: { inicio: inicio, fin: fin, base: base },
        type: 'POST',
        dataType: 'html',
        success: function(result) {
            console.log("Operacion Realizada con Exito");
            $("#maincontainer2").html(result);
        },
        error: function(xhr, status) {
            console.log("Error en la consulta")
        },
        complete: function(xhr, status) {
            console.log("Proceso Terminado")
        }
    });
}

function FuncionX() {
    $("#main2").html("HOLA");
}


function Tabla(tipo) {
    inicio = $("#inicio").val();
    fin = $("#fin").val();
    vendedor = $("#vendedor").val();

    if (inicio == "" && fin == "" && vendedor == "") {
        console.log("Todos Vacios");
        alert("Ingresa algún filtro");
        return false;
    }

    if (inicio != "" && fin == "" || inicio == "" && fin != "") {
        console.log("Al menos una fecha vacia");
        alert("Ingresa correctamente las fechas");
        return false;
    }


    if (tipo == "Viviendas" || tipo == "ViviendasTabla") {
        var1 = $("#inicio").is(":checked");
        var2 = $("#fin").is(":checked");

        if (var1 == true && var2 == false) {
            inicio = "si"
            fin = "no"
        }

        if (var1 == true && var2 == true) {
            inicio = "si"
            fin = "si"
        }

        if (var1 == false && var2 == true) {
            inicio = "no"
            fin = "si"
        }

        if (var1 == false && var2 == false && vendedor == "") {
            console.log("Todos Vacios");
            alert("Ingresa algún filtro");
            return false;
        }
    }



    $.ajax({
        url: '/tablas',
        data: { inicio: inicio, fin: fin, vendedor: vendedor, tipo: tipo },
        type: 'POST',
        dataType: 'html',
        success: function(result) {
            console.log("Operacion Realizada con Exito");
            $("#tablacontrato").html(result);
        },
        error: function(xhr, status) {
            console.log("Error en la consulta")
        },
        complete: function(xhr, status) {
            console.log("Proceso Terminado")
        }
    });
}

function GetCoworker(info) {

    $.ajax({
        url: '/workers',
        data: { info: info },
        type: 'POST',
        dataType: 'html',
        success: function(result) {
            console.log("Operacion Realizada con Exito");
            $("#maincontainer").html(result);
        },
        error: function(xhr, status) {
            console.log("Error en la consulta")
        },
        complete: function(xhr, status) {
            console.log("Proceso Terminado")
        }
    });
}

function GetCoworker2(info) {

    buscar = $("#" + info).val();
    bs = buscar.toUpperCase();
    hides = $("#hides").val();

    $.ajax({
        url: '/gettable',
        data: { bs: bs, hides: hides },
        type: 'POST',
        dataType: 'html',
        success: function(result) {
            console.log("Operacion Realizada con Exito");
            $("#maincontainer2").html(result);
        },
        error: function(xhr, status) {
            console.log("Error en la consulta")
        },
        complete: function(xhr, status) {
            console.log("Proceso Terminado")
        }
    });
}

function ImprimirArchivo(user) {

    $.ajax({
        url: '/imprimir',
        data: { user: user },
        type: 'POST',
        dataType: 'html',
        success: function(result) {
            console.log("Operacion Realizada con Exito");
            // $("#maincontainer").html(result);
            alertify.success("Documento Generado");
        },
        error: function(xhr, status) {
            console.log("Error en la consulta")
        },
        complete: function(xhr, status) {
            console.log("Proceso Terminado")
        }
    });
}


// function EnviarArchivo(){
//     var formData = new FormData($("#filer_input")[0]);

//     console.log("Llegamos");
//     console.log(formData);
//      $.ajax({
//          url: '/subirarchivo',
//          type: 'POST',
//          data: formData,
//          cache: false,
//          processData: false,
//          contentType: false,
//          enctype: 'multipart/form-data',
//          success: function (result) {
//            console.log("Regresamos");
//            Mifuncion('Prospectos');
//          },
//          error : function(xhr, status){
//              console.log("Error en la consulta")
//          },
//          complete : function(xhr, status){
//              console.log("Proceso Terminado")
//          }
//      });
// }


function AgregarColaborador() {

    $.ajax({
        url: '/userform',
        data: {},
        type: 'POST',
        dataType: 'html',
        success: function(result) {
            console.log("Operacion Realizada con Exito");
            $("#buscadorcontainer").html(result);
            $("#maincontainer").html("");
            alertify.success("Se aceptan campos vacios.");
        },
        error: function(xhr, status) {
            console.log("Error en la consulta")
        },
        complete: function(xhr, status) {
            console.log("Proceso Terminado")
        }
    });
}

function Editar(user) {

    $.ajax({
        url: '/editarcolaborador',
        data: { user: user },
        type: 'POST',
        dataType: 'html',
        success: function(result) {
            console.log("Operacion Realizada con Exito");
            $("#maincontainer").html(result);
            alertify.success("Se aceptan campos vacios.");
        },
        error: function(xhr, status) {
            console.log("Error en la consulta")
        },
        complete: function(xhr, status) {
            console.log("Proceso Terminado")
        }
    });

}

function Ver(user) {

    $.ajax({
        url: '/verpersonal',
        data: { user: user },
        type: 'POST',
        dataType: 'html',
        success: function(result) {
            console.log("Operacion Realizada con Exito");
            $("#maincontainer").html(result);
            alertify.success("Se aceptan campos vacios.");
        },
        error: function(xhr, status) {
            console.log("Error en la consulta")
        },
        complete: function(xhr, status) {
            console.log("Proceso Terminado")
        }
    });

}


function ValidaContratos() {
    value = $("#tipocontrato").val();

    if (value != "DETERMINADO") {
        $("#dias").val("INDETERMINADO");
        $("#dias").prop('readonly', true);
        $("#fincontrato").prop('disabled', true);
    } else {
        $("#dias").val("");
        $("#dias").prop('readonly', false);
        $("#fincontrato").prop('disabled', false);
        return false;
    }
}


function CambiarMunicipio() {
    value = $("#estado").val();

    if (value != "") {
        // $("#municipio").prop('readonly', false);
        $.ajax({
            url: '/getmunicipio',
            data: { value: value },
            type: 'POST',
            dataType: 'html',
            success: function(result) {
                console.log("Operacion Realizada con Exito");
                $("#municipio").html(result);
            },
            error: function(xhr, status) {
                console.log("Error en la consulta")
            },
            complete: function(xhr, status) {
                console.log("Proceso Terminado")
            }
        });

        console.log("cambiando" + value);

    }
    return false;
}

function Activos() {
    $.ajax({
        url: '/gettable',
        data: {},
        type: 'POST',
        dataType: 'html',
        success: function(result) {
            console.log("Operacion Realizada con Exito");
            $("#maincontainer").html(result);
        },
        error: function(xhr, status) {
            console.log("Error en la consulta")
        },
        complete: function(xhr, status) {
            console.log("Proceso Terminado")
        }
    });
}