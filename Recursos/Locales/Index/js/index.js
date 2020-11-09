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
            console.log("Operacion Realizada con Exito");
            document.body.style.backgroundImage = 'url("Recursos/Imagenes/Fondos/' + result + '.jpg")';
            data2 = parseInt(result) + 1;
            document.getElementById('cambiafondos').href = 'Javascript:CambiarFondo(' + data2 + ')';

            if (parseInt(result) == 33) {
                document.body.style.backgroundImage = 'url("Recursos/Imagenes/Fondos/' + 1 + '.jpg")';
                document.getElementById('cambiafondos').href = "Javascript:CambiarFondo(1)";

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
    document.body.style.color = "#ffffff";

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

    var r = confirm("Cerrar Venta de Mesa");
    fecha = $("#fechaarchivo").val();

    data += ":" + fecha

    if (r == true) {
        console.log(data);
        $.ajax({
            url: '/cerrarventademesa',
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
        VerificaMesas(fecha);

    } else {
        return false;
    }

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

$('#ModalProductos').on('hidden.bs.modal', function(e) {
    $("#modalcontainerforprod").html("");
    $("#serviciosmodal").html("");

    fecha = $("#fechaarchivo").val();
    console.log("Fecha al cerrar el modal " + fecha);

    // inputhidden = document.createElement("input");
    // inputhidden.type = "hidden";
    // inputhidden.name = "fechaarchivo";
    // inputhidden.value = fecha

    // form = document.getElementById("formventas")
    // form.appendChild(inputhidden);

    VerificaMesas(fecha);



    // form.submit();

});

function AgregarAMesa(data) {

    console.log(data);

}

function ObtenDatosDeMesa(data) {

    idmesa = $("#ModalProductos input[name=idmesa]").val();
    console.log(" idmesa", idmesa)

    data = idmesa + ":" + data
    console.log(" idmesa + data ", data);

    cadena = data.split(":");

    console.log(cadena)

    console.log(data);
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
    console.log(" idmesa", idmesa)

    data = idmesa + ":" + data
    console.log(" idmesa + data ", data);

    cadena = data.split(":");

    console.log(cadena)

    console.log(data);
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
        // conversionpromo = $("#conversionpromo").html(data);
        // conversionpromo = $("#conversionpromo");

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

        // a2 = document.createElement("a");
        // a2.classList.add("btn");
        // a2.classList.add("btn-success");
        // a2.setAttribute("role", "button");
        // a2.setAttribute("href", "Javascript:alert('a2');");
        // a2.innerHTML = "-";

        conversionpromo.appendChild(label1);
        conversionpromo.appendChild(a1);
        // conversionpromo.appendChild(a2);


        console.log(conversionpromo);

    } else {
        $("#conversionpromo").html("");
        $("#promoproductos").html("");
        return false;
    }

}

function AgregarArticulo() {

    // var conversionpromo = document.getElementById("conversionpromo");
    var promoproductos = document.getElementById("promoproductos");

    hijos = promoproductos.childElementCount;

    articulist = document.getElementsByClassName("articulito");


    numarti = articulist.length;
    var articulosagregados = document.getElementById("articulosagregados");
    articulosagregados.value = numarti + 1

    // alert("tiene " + hijos)

    divsm4 = document.createElement("div");
    divsm2 = document.createElement("div");
    divsm22 = document.createElement("div");

    divsm4.classList.add("col-sm-2");
    divsm2.classList.add("col-sm-4");
    divsm22.classList.add("col-sm-2");




    select1 = document.createElement("select");
    select1.classList.add("form-control");


    // option1 = document.createElement("option");
    // option1.setAttribute("value", "Selecciona el articulo");
    // option1.innerHTML = "opcion 1";
    // select1.appendChild(option1);

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

    // select1.innerHTML = `<option value="familiar">Familiar</option>
    // <option value="media">Coronita Media</option>`;


    select1.setAttribute("name", "productoid" + (numarti));
    select1.setAttribute("id", "productoid" + (numarti));

    // input1 = document.createElement("input");
    // input1.setAttribute("type", "");
    // input1.setAttribute("name", "productoname" + (hijos + 1));
    // input1.setAttribute("id", "poductoid" + (hijos + 1));
    // input1.setAttribute("value", "poductoid" + (hijos + 1));


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

    //i=5 i<=10 5 6 7 8 9






    // promoproductos.removeChild(promoproductos.lastChild);
    // promoproductos.removeChild(promoproductos.lastChild);
    // promoproductos.removeChild(promoproductos.lastChild);
    // promoproductos.removeChild(promoproductos.lastChild);
    // promoproductos.removeChild(promoproductos.lastChild);
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