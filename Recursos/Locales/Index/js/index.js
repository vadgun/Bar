let indexit = document.getElementById("indexit")
let administracionit = document.getElementById("administracionit")
let inventariosit = document.getElementById("inventariosit")
let ventasit = document.getElementById("ventasit")
let ventasdiariasit = document.getElementById("ventasdiariasit")
let imprimirventasdiariasit = document.getElementById("imprimirventasdiariasit")

$(document).ready(function() {
    //Saber el ultimo fondo puesto y atribuirlo al body.style
    $.ajax({
        url: '/ultimofondousado',
        data: {},
        type: 'POST',
        dataType: 'html',
        success: function(result) {
            console.log("Operacion Realizada con Exito FONDO");
            var parseado = JSON.parse(result);
            document.body.style.backgroundImage = 'url("Recursos/Imagenes/Fondos/' + parseado.Disponibles + '.jpg")';
            data2 = parseInt(parseado.Disponibles) + 1;
            document.getElementById('cambiafondos').href = 'Javascript:CambiarFondo(' + data2 + ')';

            if (parseInt(result) == 33) {
                document.body.style.backgroundImage = 'url("Recursos/Imagenes/Fondos/' + 1 + '.jpg")';
                document.getElementById('cambiafondos').href = "Javascript:CambiarFondo(1)";
            }

            existes = document.getElementById('message1');

            if (existes != undefined){
                document.getElementById('message1').innerHTML = parseado.Mensajes[0];
                document.getElementById('message2').innerHTML = parseado.Mensajes[1];
                document.getElementById('message3').innerHTML = parseado.Mensajes[2];
            }    
        },
        error: function(xhr, status) {
            console.log("Error en la consulta")
        },
        complete: function(xhr, status) {
            console.log("Proceso Terminado")
        }
    });
    
});

function CambiarFondo(data) {

    document.body.style.backgroundImage = 'url("Recursos/Imagenes/Fondos/' + data + '.jpg")';
    data2 = data + 1;
    document.getElementById('cambiafondos').href = 'Javascript:CambiarFondo(' + data2 + ')';

    if (data == 33) {
        document.body.style.backgroundImage = 'url("Recursos/Imagenes/Fondos/' + 1 + '.jpg")';
        document.getElementById('cambiafondos').href = "Javascript:CambiarFondo(1)";

    }

    $.ajax({
        url: '/cambiarfondo',
        data: { data: data },
        type: 'POST',
        dataType: 'html',
        success: function(result) {
            console.log("Operacion Realizada con Exito");
        },
        error: function(xhr, status) {
            console.log("Error en la consulta")
        },
        complete: function(xhr, status) {
            console.log("Proceso Terminado")
        }
    });

}

function EliminarVenta(data) {
    console.log(data);
    $.ajax({
        url: '/vaciarbasededatos',
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

function LetrasBlancas() {
    document.body.style.color = "#99bbff";

}

function LetrasNegras() {
    document.body.style.color = "#000000";
}


function Mifuncion(data) {
    console.log(data);
    $.ajax({
        url: '/altaform',
        data: { data: data },
        type: 'POST',
        dataType: 'html',
        success: function(result) {
            console.log("Operacion Realizada con Exito");
            $("#maincontainer").html(result);
            $("#subcontainer").html("");
        },
        error: function(xhr, status) {
            console.log("Error en la consulta")
        },
        complete: function(xhr, status) {
            console.log("Proceso Terminado")
        }
    });
}


function EditarAlmacen(data) {
    console.log(data);
    $.ajax({
        url: '/editaralmacen',
        data: { data: data },
        type: 'POST',
        dataType: 'html',
        success: function(result) {
            console.log("Operacion Realizada con Exito");
            $("#subcontainer").html(result);
        },
        error: function(xhr, status) {
            console.log("Error en la consulta")
        },
        complete: function(xhr, status) {
            console.log("Proceso Terminado2")
        }
    });
}

function ClasificaVenta(data) {
    console.log(data);
    $.ajax({
        url: '/clasificaventa',
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
            console.log("Proceso Terminado2")
        }
    });


}

function ImprimeVenta(data) {
    console.log(data);
    $.ajax({
        url: '/imprimeventa',
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
            console.log("Proceso Imprimir Terminado")
        }
    });


}

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



function VerificaMesas(data) {
    console.log(data);
    $.ajax({
        url: '/verificarventa',
        data: { data: data },
        type: 'POST',
        dataType: 'html',
        success: function(result) {
            console.log("Operacion Realizada con Exito");

            $("#maincontainer").html(result);
        },
        error: function(xhr, status) {
            console.log("Error en la consulta");
        },
        complete: function(xhr, status) {
            console.log("Proceso Verificar Mesa Terminado");
        }
    });
}


function VerificaVentaDiaria(data) {
    console.log(data);
    $.ajax({
        url: '/verificarventadiaria',
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
            console.log("Proceso Verificar Venta Diaria Terminado");
        }
    });
}

function AbrirVentaDelDia(data) {
    console.log(data);
    $.ajax({
        url: '/abrirventadiaria',
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
            console.log("Proceso Terminado2")
        }
    });
}

//Eliminar el producto si no lo encuentra en existencia en los alamacenes
function EliminarProducto(data) {

    var r = confirm("Eliminar Producto");
    if (r == true) {
        console.log(data);
        $.ajax({
            url: '/eliminarproducto',
            data: { data: data },
            type: 'POST',
            dataType: 'html',
            success: function(result) {
                console.log("Operacion Realizada con Exito");
                $("#subsubcontainer").html(result);
            },
            error: function(xhr, status) {
                console.log("Error en la consulta")
            },
            complete: function(xhr, status) {
                console.log("Proceso Terminado2")
            }
        });
    } else {
        return false;
    }


}

//Edicion de Productos
function EditarProducto(data) {
    console.log(data);
    $.ajax({
        url: '/editarproducto',
        data: { data: data },
        type: 'POST',
        dataType: 'html',
        success: function(result) {
            console.log("Operacion Realizada con Exito");
            $("#subsubcontainer").html(result);
        },
        error: function(xhr, status) {
            console.log("Error en la consulta")
        },
        complete: function(xhr, status) {
            console.log("Proceso Terminado2")
        }
    });
}

//Funciones de Traspasos 
function EnviarBodega(data) {
    console.log(data);
    $.ajax({
        url: '/enviarproductoabodega',
        data: { data: data },
        type: 'POST',
        dataType: 'html',
        success: function(result) {
            console.log("Operacion Realizada con Exito");
            $("#subsubcontainer").html(result);
        },
        error: function(xhr, status) {
            console.log("Error en la consulta")
        },
        complete: function(xhr, status) {
            console.log("Proceso Terminado2")
        }
    });
}

function EnviarRefrigerador(data) {
    console.log(data);
    $.ajax({
        url: '/enviarproductoarefrigerador',
        data: { data: data },
        type: 'POST',
        dataType: 'html',
        success: function(result) {
            console.log("Operacion Realizada con Exito");
            $("#subsubcontainer").html(result);
        },
        error: function(xhr, status) {
            console.log("Error en la consulta")
        },
        complete: function(xhr, status) {
            console.log("Proceso Terminado2")
        }
    });
}


function EliminarDeBodega(data) {
    console.log(data);
    $.ajax({
        url: '/eliminarproductodebodega',
        data: { data: data },
        type: 'POST',
        dataType: 'html',
        success: function(result) {
            console.log("Operacion Realizada con Exito");
            $("#subsubcontainer").html(result);
        },
        error: function(xhr, status) {
            console.log("Error en la consulta")
        },
        complete: function(xhr, status) {
            console.log("Proceso Terminado2")
        }
    });
}

function EliminarDeRefrigerador(data) {
    console.log(data);
    $.ajax({
        url: '/eliminarproductoderefrigerador',
        data: { data: data },
        type: 'POST',
        dataType: 'html',
        success: function(result) {
            console.log("Operacion Realizada con Exito");
            $("#subsubcontainer").html(result);
        },
        error: function(xhr, status) {
            console.log("Error en la consulta")
        },
        complete: function(xhr, status) {
            console.log("Proceso Terminado2")
        }
    });
}

function CerrarVenta(data) {
    const buttons = document.querySelectorAll('[data-mesa]');
    var numesa = ""
    
    for (x of buttons) {
        let id = x.getAttribute('data-id');
        if (id == data){
            numesa = x.getAttribute('data-mesa');
        }
    }

    fecha = $("#fechaarchivo").val();
    data += ":" + fecha

    alertify.confirm("Desea cerrar la mesa "+numesa+" ?",
    function(){
        alertify.success('Ok');
                $.ajax({
            url: '/cerrarventademesa',
            data: { data: data },
            type: 'POST',
            dataType: 'html',
            success: function(result) {
                console.log("Operacion Realizada con Exito");
                console.log(result);
            },
            error: function(xhr, status) {
                console.log("Error en la consulta")
            },
            complete: function(xhr, status) {
                console.log("Proceso Terminado2")
            }
        });
    },
    function(){
        alertify.error('Cancel');
    });

}

$('#ModalProductos').on('show.bs.modal', function(e) {
    var mesa = $(e.relatedTarget).data('mesa');
    var id = $(e.relatedTarget).data('id');
    $("#ModalProductos input[name=idmesa]").val(id);
    $("#titulodemesa").html("Mesa en venta : " + mesa);
    console.log(" mesa", mesa, "   -  id    --", id)
    data = id;
    $.ajax({
        url: '/productosVendidos',
        data: { data: data },
        type: 'POST',
        dataType: 'html',
        success: function(result) {
            console.log("Operacion Realizada con Exito");
            $("#serviciosmodal").html(result);
        },
        error: function(xhr, status) {
            console.log("Error en la consulta")
        },
        complete: function(xhr, status) {
            console.log("Datos de Mesa obtenidos")

        }
    });
});

function ProductosYaVendidos(data){
    $.ajax({
        url: '/productosVendidos',
        data: { data: data },
        type: 'POST',
        dataType: 'html',
        success: function(result) {
            console.log("Operacion Realizada con Exito");
            $("#serviciosmodal").html(result);
        },
        error: function(xhr, status) {
            console.log("Error en la consulta")
        },
        complete: function(xhr, status) {
            console.log("Datos de Mesa obtenidos")

        }
    });
}

$('#ModalProductos').on('hidden.bs.modal', function(e) {
    $("#modalcontainerforprod").html("");
    $("#serviciosmodal").html("");
});

function AgregarAMesa(data) {
    console.log(data);
}

function ObtenDatosDeMesa(data) {

    idmesa = $("#ModalProductos input[name=idmesa]").val();
    data = idmesa + ":" + data
    cadena = data.split(":");
    console.log(cadena)
    $.ajax({
        url: '/productosAgregadosEnModal',
        data: { data: data },
        type: 'POST',
        dataType: 'html',
        success: function(result) {
            console.log("Operacion Realizada con Exito");
            $("#serviciosmodal").html(result);
        },
        error: function(xhr, status) {
            console.log("Error en la consulta")
        },
        complete: function(xhr, status) {
            console.log("Datos de Mesa obtenidos")
            MifuncionModal(cadena[3]);
        }
    });
}

function ObtenDatosDeMesaDesdePromocion(data) {

    idmesa = $("#ModalProductos input[name=idmesa]").val();
    data = idmesa + ":" + data
    cadena = data.split(":");
    $.ajax({
        url: '/productosAgregadosEnModalDesdePromo',
        data: { data: data },
        type: 'POST',
        dataType: 'html',
        success: function(result) {
            console.log("Operacion Realizada con Exito");
            $("#serviciosmodal").html(result);
        },
        error: function(xhr, status) {
            console.log("Error en la consulta")
        },
        complete: function(xhr, status) {
            console.log("Datos de Mesa obtenidos")
            MifuncionModal(cadena[3]);
        }
    });

}

function ConvierteAPromo(data) {

    if (data == "promo") {
        var conversionpromo = document.getElementById("conversionpromo");

        label1 = document.createElement("label");
        label1.classList.add("col-sm-2");
        label1.classList.add("col-form-label");
        label1.classList.add("negrita");
        label1.innerHTML = "Agrega Productos:";


        a1 = document.createElement("a");
        a1.classList.add("btn");
        a1.classList.add("btn-success");
        a1.setAttribute("role", "button");
        a1.innerHTML = "+";
        a1.setAttribute("href", "Javascript:AgregarArticulo();");

        conversionpromo.appendChild(label1);
        conversionpromo.appendChild(a1);
    } else {
        $("#conversionpromo").html("");
        $("#promoproductos").html("");
        return false;
    }

}

function AgregarArticulo() {
    var promoproductos = document.getElementById("promoproductos");
    hijos = promoproductos.childElementCount;

    articulist = document.getElementsByClassName("articulito");
    numarti = articulist.length;

    var articulosagregados = document.getElementById("articulosagregados");
    articulosagregados.value = numarti + 1

    divsm4 = document.createElement("div");
    divsm2 = document.createElement("div");
    divsm22 = document.createElement("div");

    divsm4.classList.add("col-sm-2");
    divsm2.classList.add("col-sm-4");
    divsm22.classList.add("col-sm-2");

    select1 = document.createElement("select");
    select1.classList.add("form-control");

    $.ajax({
        url: '/productosASelect',
        data: {},
        type: 'POST',
        dataType: 'html',
        success: function(result) {
            console.log("Operacion Realizada con Exito");
            select1.innerHTML = result
        },
        error: function(xhr, status) {
            console.log("Error en la consulta")
        },
        complete: function(xhr, status) {
            console.log("Productos insertados en el select")
        }
    });

    select1.setAttribute("name", "productoid" + (numarti));
    select1.setAttribute("id", "productoid" + (numarti));

    input2 = document.createElement("input");
    input2.setAttribute("type", "number");
    input2.setAttribute("name", "cantidadprod" + (numarti));
    input2.setAttribute("id", "cantidadpro" + (numarti));
    input2.classList.add("form-control");

    label1 = document.createElement("label");
    label1.classList.add("col-sm-2");
    label1.classList.add("col-form-label");
    label1.classList.add("negrita");
    label1.classList.add("articulito");
    label1.innerHTML = "Articulo " + (numarti + 1) + " :";

    divsm2.appendChild(select1);
    divsm4.appendChild(input2);

    ax = document.createElement("a");
    ax.classList.add("btn");
    ax.classList.add("btn-danger");
    ax.setAttribute("role", "button");
    ax.innerHTML = "Quitar";
    ax.setAttribute("href", "Javascript:EliminarArticulo();");


    promoproductos.appendChild(label1);
    promoproductos.appendChild(divsm2);
    promoproductos.appendChild(divsm4);
    promoproductos.appendChild(divsm22);
    // promoproductos.appendChild(select1);
    // promoproductos.appendChild(input2);
    promoproductos.appendChild(ax);

}

function EliminarArticulo() {
    var promoproductos = document.getElementById("promoproductos");

    hijos = promoproductos.childElementCount;

    console.log(hijos);
    //i=10  i<10  No

    for (i = hijos; i > (hijos - 5); i--) {
        console.log(i - 1);
        promoproductos.removeChild(promoproductos.childNodes[i]);
    }

    articulist = document.getElementsByClassName("articulito");
    numarti = articulist.length;
    var articulosagregados = document.getElementById("articulosagregados");
    articulosagregados.value = numarti;
}


function MifuncionModal(data) {

    console.log(data);

    $.ajax({
        url: '/productosEnModal',
        data: { data: data },
        type: 'POST',
        dataType: 'html',
        success: function(result) {
            console.log("Operacion Realizada con Exito");
            $("#modalcontainerforprod").html(result);
        },
        error: function(xhr, status) {
            console.log("Error en la consulta")
        },
        complete: function(xhr, status) {
            console.log("Productos Cargados para Agregar")

        }
    });

}

function CrearTabla() {

    var formulario = document.createElement("form");
    document.getElementById("tablecontainer").appendChild(formulario);

    // $("#tablecontainer").html(result);
}