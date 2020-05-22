var HISTORIAL_ID = -1;
var IDENTIFICACION = "";

function busquedaDNI(event){
    if(document.querySelector("#inputDNI").value.length != 9){
        //Activar alerta
        document.querySelector("#alert").textContent = "El documento de identificación debe tener un formato válido (por ejemplo, 00000000X)";
        document.querySelector("#alert").classList.replace("alert-success", "alert-danger");
        document.querySelector("#alert").classList.remove('invisible');
        document.querySelector("#historialDiv").classList.add('invisible');
        return;
    }else{
        restBuscarDNI(document.querySelector("#inputDNI").value);
        IDENTIFICACION = document.querySelector("#inputDNI").value;
    }
}

function paramBusquedaDNI(identificacion){
    if(identificacion.length != 9){
        //Activar alerta
        document.querySelector("#alert").textContent = "El documento de identificación debe tener un formato válido (por ejemplo, 00000000X)";
        document.querySelector("#alert").classList.replace("alert-success", "alert-danger");
        document.querySelector("#alert").classList.remove('invisible');
        document.querySelector("#historialDiv").classList.add('invisible');
        return;
    }else{
        restBuscarDNI(identificacion);
        IDENTIFICACION = identificacion;
    }
}

function cargarTablaHistorial(eList){
    //Eliminamos las filas de la tabla
    while (document.querySelector("#historialTablaBody").lastElementChild) { 
        document.querySelector("#historialTablaBody").removeChild(document.querySelector("#historialTablaBody").lastElementChild); 
    } 
    if(!eList || eList.length < 1){
        document.querySelector("#historialTabla").classList.add('invisible');
    }else{
        document.querySelector("#alert").classList.add('invisible');
        document.querySelector("#historialTabla").classList.remove('invisible');
        eList.forEach(entrada => {
            addRowHistorial(entrada);
        });
    }
}

function cargarTablaAnaliticas(aList){
    //Eliminamos las filas de la tabla
    while (document.querySelector("#analiticasTablaBody").lastElementChild) { 
        document.querySelector("#analiticasTablaBody").removeChild(document.querySelector("#analiticasTablaBody").lastElementChild); 
    } 
    if(!aList || aList.length < 1){
        document.querySelector("#analiticasTabla").classList.add('invisible');
    }else{
        document.querySelector("#alert").classList.add('invisible');
        document.querySelector("#analiticasTabla").classList.remove('invisible');
        aList.forEach(analitica => {
            addRowAnaliticas(analitica);
        });
    }
}

function consultarEntradaHistorial(event){
    window.location.href = "/user/emergency/historial/entrada?entradaId=" + event.target.closest("tr").getAttribute("id") + "&identificacion="+IDENTIFICACION;
}

function consultarAnaliticaHistorial(event){
    window.location.href = "/user/emergency/historial/analitica?analiticaId=" + event.target.closest("tr").getAttribute("id") + "&identificacion="+IDENTIFICACION;
}

function addRowHistorial(entrada){
    let tr = document.createElement('tr');
    let fecha = document.createElement('td');
    let especialista = document.createElement('td');
    let tipo = document.createElement('td');
    let acciones = document.createElement('td');

    let accionesButton = document.createElement('button');
    accionesButton.classList = "btn btn-primary";
    accionesButton.type = "button";
    accionesButton.textContent = "Consultar entrada";
    accionesButton.addEventListener("click", consultarEntradaHistorial, false);
    acciones.append(accionesButton);
    fecha.textContent = entrada.createdAt;
    especialista.textContent = entrada.empleadoNombre;
    tipo.textContent = entrada.tipo;
    tr.append(tipo);
    tr.append(especialista);
    tr.append(fecha);
    tr.append(acciones);
    tr.setAttribute("id", entrada.id);
    //Añadimos fila a la tabla
    document.querySelector(`#historialTabla`).querySelector('tbody').append(tr);
}

function addRowAnaliticas(entrada){
    let tr = document.createElement('tr');
    let fecha = document.createElement('td');
    let especialista = document.createElement('td');
    let tipo = document.createElement('td');
    let acciones = document.createElement('td');

    let accionesButton = document.createElement('button');
    accionesButton.classList = "btn btn-primary";
    accionesButton.type = "button";
    accionesButton.textContent = "Consultar analítica";
    accionesButton.addEventListener("click", consultarAnaliticaHistorial, false);
    acciones.append(accionesButton);
    fecha.textContent = entrada.createdAt;
    especialista.textContent = entrada.empleadoNombre;
    tipo.textContent = "Analítica";
    tr.append(tipo);
    tr.append(especialista);
    tr.append(fecha);
    tr.append(acciones);
    tr.setAttribute("id", entrada.id);
    //Añadimos fila a la tabla
    document.querySelector(`#analiticasTabla`).querySelector('tbody').append(tr);
}

function restBuscarDNI(DNI){
    const url= `/user/emergency/historial`;
    const payload= {identificacion: DNI};
    const request = {
        method: 'POST', 
        headers: cabeceras,
        body: JSON.stringify(payload),
    };
    fetch(url,request)
    .then( response => response.json() )
        .then( r => {
            if(!r.Error){
                if(r.id != 0){
                    //PROCESAR HISTORIAL
                    HISTORIAL_ID = r.id;
                    document.querySelector("#alert").classList.add('invisible');
                    document.querySelector("#historialDiv").classList.remove('invisible');
                    document.querySelector("#spanNombre").textContent = r.nombrePaciente + " " + r.apellidosPaciente;
                    document.querySelector("#spanSexo").textContent = r.sexo;
                    document.querySelector("#spanAlergias").textContent = r.alergias
                    if(r.entradas == null){
                        document.querySelector("#alertTablaHistorial").classList.remove('invisible');
                    }else{
                        document.querySelector("#alertTablaHistorial").classList.add('invisible');
                    }
                    if(r.analiticas == null){
                        document.querySelector("#alertTablaAnaliticas").classList.remove('invisible');
                    }else{
                        document.querySelector("#alertTablaAnaliticas").classList.add('invisible');
                    }
                    cargarTablaHistorial(r.entradas);
                    cargarTablaAnaliticas(r.analiticas);
                    
                }else{
                    document.querySelector("#historialDiv").classList.add('invisible');
                    document.querySelector("#alert").textContent = "No existe ningún usuario con esa identificación";
                    document.querySelector("#alert").classList.replace("alert-success", "alert-danger");
                    document.querySelector("#alert").classList.remove('invisible');
                }
            }
            else{
                document.querySelector("#alert").textContent = r.Error;
                document.querySelector("#alert").classList.replace("alert-success", "alert-danger");
                document.querySelector("#alert").classList.remove('invisible');
                document.querySelector("#historialDiv").classList.add('invisible');
            }
        })
        .catch(err => alert(err));
}

function addEntrada(){
    window.location.href = "/user/emergency/historial/addEntrada?historialId=" + HISTORIAL_ID + "&identificacion="+IDENTIFICACION;
}

function addAnalitica(){
    window.location.href = "/user/emergency/historial/addAnalitica?historialId=" + HISTORIAL_ID + "&identificacion="+IDENTIFICACION;
}


function init () {
    //Si se pasa por parametro el DNI se busca auto
    var url = new URL(window.location.href);
    var paramIdentificacion = url.searchParams.get("identificacion");
    if(paramIdentificacion){
        paramBusquedaDNI(paramIdentificacion);
    }
    deleteBreadcrumb();
    addLinkBreadcrumb('Usuario', '/user/menu');
    addLinkBreadcrumb('Emergencias', '/user/emergency');
    document.querySelector("#searchButton").addEventListener('click',busquedaDNI,false);
    document.querySelector("#addEntradaButton").addEventListener('click',addEntrada,false);
    document.querySelector("#addAnaliticaButton").addEventListener('click',addAnalitica,false);
}

document.addEventListener('DOMContentLoaded',init,false);