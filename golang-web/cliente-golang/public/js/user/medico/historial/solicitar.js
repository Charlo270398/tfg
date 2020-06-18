var HISTORIAL_ID = -1;
var IDENTIFICACION = "";

function busquedaDNI(event){
    //Eliminamos datos anteriores
    while (document.querySelector("#buttonsForm").lastElementChild) { 
        document.querySelector("#buttonsForm").removeChild(document.querySelector("#buttonsForm").lastElementChild); 
    } 
    while (document.querySelector("#historialTablaBody").lastElementChild) { 
        document.querySelector("#historialTablaBody").removeChild(document.querySelector("#historialTablaBody").lastElementChild); 
    } 
    while (document.querySelector("#analiticasTablaBody").lastElementChild) { 
        document.querySelector("#analiticasTablaBody").removeChild(document.querySelector("#analiticasTablaBody").lastElementChild); 
    } 

    if(document.querySelector("#inputDNI").value.length != 9){
        //Activar alerta
        document.querySelector("#alert").textContent = "El documento de identificación debe tener un formato válido (por ejemplo, 00000000X)";
        document.querySelector("#alert").classList.replace("alert-success", "alert-danger");
        document.querySelector("#alert").classList.remove('invisible');
        document.querySelector("#buttonsForm").classList.add('invisible');
        document.querySelector("#alertBusqueda").classList.add('invisible');
        document.querySelector("#historialTabla").classList.add('invisible');
        document.querySelector("#historialTablaTitulo").classList.add('invisible');
        document.querySelector("#analiticasTabla").classList.add('invisible');
        document.querySelector("#analiticasTablaTitulo").classList.add('invisible');
        document.querySelector("#datosBasicos").classList.add('invisible');
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
        document.querySelector("#historialTabla").classList.add('invisible');
        document.querySelector("#historialTablaTitulo").classList.add('invisible');
        document.querySelector("#analiticasTabla").classList.add('invisible');
        document.querySelector("#analiticasTablaTitulo").classList.add('invisible');
        document.querySelector("#datosBasicos").classList.add('invisible');
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
        document.querySelector("#historialTablaTitulo").classList.add('invisible');
    }else{
        document.querySelector("#alert").classList.add('invisible');
        document.querySelector("#historialTablaTitulo").classList.remove('invisible');
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
        document.querySelector("#analiticasTablaTitulo").classList.add('invisible');
    }else{
        document.querySelector("#alert").classList.add('invisible');
        document.querySelector("#analiticasTabla").classList.remove('invisible');
        document.querySelector("#analiticasTablaTitulo").classList.remove('invisible');
        aList.forEach(analitica => {
            addRowAnaliticas(analitica);
        });
    }
}

function restBuscarDNI(DNI){
    document.querySelector("#datosBasicos").classList.add('invisible');
    const url= `/user/doctor/historial/solicitar`;
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
                if(r.historial.id == 0){
                    //Si no existe un paciente 
                    document.querySelector("#alert").textContent = "No existe un paciente con ese nº de identificación";
                    document.querySelector("#alert").classList.replace("alert-success", "alert-danger");
                    document.querySelector("#alert").classList.remove('invisible');
                    document.querySelector("#historialTabla").classList.add('invisible');
                    document.querySelector("#historialTablaTitulo").classList.add('invisible');
                    document.querySelector("#analiticasTabla").classList.add('invisible');
                    document.querySelector("#analiticasTablaTitulo").classList.add('invisible');
                    document.querySelector("#datosBasicos").classList.add('invisible');
                    document.querySelector("#buttonsForm").classList.add('invisible');
                    document.querySelector("#alertBusqueda").classList.add('invisible');
                }
                else{
                    //PROCESAR HISTORIAL
                    HISTORIAL_ID = r.historial.id;
                    document.querySelector("#alertBusqueda").classList.add('invisible');
                    document.querySelector("#buttonsForm").classList.remove('invisible');
                    document.querySelector("#alert").classList.add('invisible');
                    procesarHistorial(r);
                }
            }
            else{
                document.querySelector("#alert").textContent = r.Error;
                document.querySelector("#alert").classList.replace("alert-success", "alert-danger");
                document.querySelector("#alert").classList.remove('invisible');
                document.querySelector("#historialTabla").classList.add('invisible');
                document.querySelector("#buttonsForm").classList.add('invisible');
                document.querySelector("#alertBusqueda").classList.add('invisible');
            }
        })
        .catch(err => alert(err));
}

function addEntrada(){
    window.location.href = "/user/doctor/historial/addEntrada?historialId=" + HISTORIAL_ID + "&identificacion="+IDENTIFICACION;
}

function addAnalitica(){
    window.location.href = "/user/doctor/historial/addAnalitica?historialId=" + HISTORIAL_ID + "&identificacion="+IDENTIFICACION;
}

function solicitarAccesoTotal(event){
    //Enviamos peticion
    const url= `/permisos/historial/total/solicitar`;
    const payload= {id: HISTORIAL_ID};
    const request = {
        method: 'POST', 
        headers: cabeceras,
        body: JSON.stringify(payload),
    };
    fetch(url,request)
    .then( response => response.json() )
        .then( r => {
            if(!r.Error){
                if(r.result == "OK"){
                    console.log("PETICION REALIZADA CORRECTAMENTE");
                }
            }
        })
        .catch(err => alert(err));

    //Eliminamos boton
    while (document.querySelector("#buttonsForm").lastElementChild) { 
        document.querySelector("#buttonsForm").removeChild(document.querySelector("#buttonsForm").lastElementChild); 
    } 
    let solicitado = document.createElement('span');
    solicitado.textContent = "Permiso total solicitado";
    document.querySelector("#buttonsForm").append(solicitado);
}

function solicitarAccesoBasico(event){
    //Enviamos peticion
    const url= `/permisos/historial/basico/solicitar`;
    const payload= {id: HISTORIAL_ID};
    const request = {
        method: 'POST', 
        headers: cabeceras,
        body: JSON.stringify(payload),
    };
    fetch(url,request)
    .then( response => response.json() )
        .then( r => {
            if(!r.Error){
                if(r.result == "OK"){
                    console.log("PETICION REALIZADA CORRECTAMENTE");
                }
            }
        })
    .catch(err => alert(err));

    //Eliminamos boton
    while (document.querySelector("#buttonsForm").lastElementChild) { 
        document.querySelector("#buttonsForm").removeChild(document.querySelector("#buttonsForm").lastElementChild); 
    } 
    let solicitado = document.createElement('span');
    solicitado.textContent = "Permiso básico solicitado";
    document.querySelector("#buttonsForm").append(solicitado);
}

function solicitarAccesoEntrada(event){
    entradaId = event.target.closest("tr").getAttribute("id");
    //Enviamos peticion
    const url= `/permisos/entrada/solicitar`;
    const payload= {id: parseInt(entradaId)};
    const request = {
        method: 'POST', 
        headers: cabeceras,
        body: JSON.stringify(payload),
    };
    fetch(url,request)
    .then( response => response.json() )
        .then( r => {
            if(!r.Error){
                if(r.result == "OK"){
                    console.log("PETICION REALIZADA CORRECTAMENTE");
                }
            }
        })
    .catch(err => alert(err));
    //Eliminamos boton
    let solicitado = document.createElement('span');
    solicitado.textContent = "Permiso solicitado";
    event.target.parentNode.append(solicitado);
    event.target.remove();
}

function solicitarAccesoAnalítica(){
    analiticaId = event.target.closest("tr").getAttribute("id");
    //Enviamos peticion
    const url= `/permisos/analitica/solicitar`;
    const payload= {id: parseInt(analiticaId)};
    const request = {
        method: 'POST', 
        headers: cabeceras,
        body: JSON.stringify(payload),
    };
    fetch(url,request)
    .then( response => response.json() )
        .then( r => {
            if(!r.Error){
                if(r.result == "OK"){
                    console.log("PETICION REALIZADA CORRECTAMENTE");
                }
            }
        })
    .catch(err => alert(err));
    //Eliminamos boton
    let solicitado = document.createElement('span');
    solicitado.textContent = "Permiso solicitado";
    event.target.parentNode.append(solicitado);
    event.target.remove();
}

function procesarHistorial(solicitud){
    //Boton añadir entrada
    let addEntradaButton = document.createElement('button');
    addEntradaButton.classList = "btn btn-primary";
    addEntradaButton.type = "button";
    addEntradaButton.textContent = "Añadir entrada";
    addEntradaButton.addEventListener("click", addEntrada, false);
    document.querySelector("#buttonsForm").append(addEntradaButton);
    //Boton añadir analítica
    let addAnaliticaButton = document.createElement('button');
    addAnaliticaButton.classList = "btn btn-primary";
    addAnaliticaButton.type = "button";
    addAnaliticaButton.textContent = "Añadir analítica";
    addAnaliticaButton.addEventListener("click", addAnalitica, false);
    document.querySelector("#buttonsForm").append(addAnaliticaButton);
    //Boton acceso total
    let accesoTotalButton = document.createElement('button');
    accesoTotalButton.classList = "btn btn-primary";
    accesoTotalButton.type = "button";
    accesoTotalButton.textContent = "Solicitar acceso total al historial";
    accesoTotalButton.addEventListener("click", solicitarAccesoTotal, false);
    document.querySelector("#buttonsForm").append(accesoTotalButton);

    //DATOS BÁSICOS
    if(solicitud.historialPermitido.id == 0){
        //No disponemos de permisos básicos
        document.querySelector("#historialTabla").classList.add('invisible');
        document.querySelector("#analiticasTabla").classList.add('invisible');
        document.querySelector("#historialTablaTitulo").classList.add('invisible');
        document.querySelector("#analiticasTablaTitulo").classList.add('invisible');
        document.querySelector("#alertBusqueda").classList.remove('invisible');
        let accesoBasicoButton = document.createElement('button');
        accesoBasicoButton.classList = "btn btn-secondary";
        accesoBasicoButton.type = "button";
        accesoBasicoButton.textContent = "Solicitar acceso a datos básicos del historial";
        accesoBasicoButton.addEventListener("click", solicitarAccesoBasico, false);
        document.querySelector("#buttonsForm").append(accesoBasicoButton);
    }else{
        //Si disponemos de permisos básicos
        document.querySelector("#datosBasicos").classList.remove('invisible');
        document.querySelector("#spanNombre").textContent = solicitud.historialPermitido.nombrePaciente + " " + solicitud.historialPermitido.apellidosPaciente;
        document.querySelector("#spanSexo").textContent = solicitud.historialPermitido.sexo;
        document.querySelector("#spanAlergias").textContent = solicitud.historialPermitido.alergias;
        //Cargamos tablas 
        cargarTablaHistorial(solicitud.historialPermitido.entradas);
        cargarTablaAnaliticas(solicitud.historialPermitido.analiticas);
    }
}

function consultarEntradaHistorial(event){
    window.location.href = "/user/medico/historial/entrada?entradaId=" + event.target.closest("tr").getAttribute("id") + "&identificacion="+IDENTIFICACION;
}

function consultarAnaliticaHistorial(event){
    window.location.href = "/user/medico/historial/analitica?analiticaId=" + event.target.closest("tr").getAttribute("id") + "&identificacion="+IDENTIFICACION;
}

function addRowHistorial(entrada){
    let tr = document.createElement('tr');
    let fecha = document.createElement('td');
    let especialista = document.createElement('td');
    let tipo = document.createElement('td');
    let acciones = document.createElement('td');

    fecha.textContent = entrada.createdAt;
    especialista.textContent = entrada.empleadoNombre;
    if(entrada.tipo != ""){
        tipo.textContent = entrada.tipo;
        let consultarEntradaButton = document.createElement('button');
        consultarEntradaButton.classList = "btn btn-primary";
        consultarEntradaButton.type = "button";
        consultarEntradaButton.textContent = "Consultar entrada";
        consultarEntradaButton.addEventListener("click", consultarEntradaHistorial, false);
        acciones.append(consultarEntradaButton);
    }
    else{
        tipo.textContent = "*Desconocido*";
        let solicitarButton = document.createElement('button');
        solicitarButton.classList = "btn btn-success";
        solicitarButton.type = "button";
        solicitarButton.textContent = "Solicitar acceso";
        solicitarButton.addEventListener("click", solicitarAccesoEntrada, false);
        acciones.append(solicitarButton);
    }
    tr.append(tipo);
    tr.append(especialista);
    tr.append(fecha);
    tr.append(acciones);
    tr.setAttribute("id", entrada.id);
    //Añadimos fila a la tabla
    document.querySelector(`#historialTabla`).querySelector('tbody').append(tr);
}

function addRowAnaliticas(analitica){
    let tr = document.createElement('tr');
    let fecha = document.createElement('td');
    let especialista = document.createElement('td');
    let tipo = document.createElement('td');
    let acciones = document.createElement('td');

    fecha.textContent = analitica.createdAt;
    especialista.textContent = analitica.empleadoNombre;
    tipo.textContent = "Analítica";
    if(analitica.clave != ""){
        let consultarAnaliticaButton = document.createElement('button');
        consultarAnaliticaButton.classList = "btn btn-primary";
        consultarAnaliticaButton.type = "button";
        consultarAnaliticaButton.textContent = "Consultar entrada";
        consultarAnaliticaButton.addEventListener("click", consultarAnaliticaHistorial, false);
        acciones.append(consultarAnaliticaButton);
    }
    else{
        let solicitarButton = document.createElement('button');
        solicitarButton.classList = "btn btn-success";
        solicitarButton.type = "button";
        solicitarButton.textContent = "Solicitar acceso";
        solicitarButton.addEventListener("click", solicitarAccesoAnalítica, false);
        acciones.append(solicitarButton);
    }
    tr.append(tipo);
    tr.append(especialista);
    tr.append(fecha);
    tr.append(acciones);
    tr.setAttribute("id", analitica.id);
    //Añadimos fila a la tabla
    document.querySelector(`#analiticasTabla`).querySelector('tbody').append(tr);
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
    addLinkBreadcrumb('Medico', '/user/doctor');
    addLinkBreadcrumb('Solicitar historial', '/user/doctor/historial/solicitar');
    document.querySelector("#searchButton").addEventListener('click',busquedaDNI,false);
}

document.addEventListener('DOMContentLoaded',init,false);