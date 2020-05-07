function init () {
    deleteBreadcrumb();
    addLinkBreadcrumb('Usuario', '/user/menu');
    addLinkBreadcrumb('Paciente', '/user/patient');
    addLinkBreadcrumb('Historia clínica', '/user/patient/historial');
    console.log(historial);
    loadTable(historial.entradas)
}

function loadTable(eList){
    if(!eList || eList.length < 1){
        document.querySelector("#alert").textContent = "No hay ninguna entrada en tu historial";
        document.querySelector("#alert").classList.replace("alert-success", "alert-danger");
        document.querySelector("#alert").classList.remove('invisible');
    }else{
        eList.forEach(entrada => {
            addRow(entrada);
        });
    }
}

function addRow(entrada){
    let tr = document.createElement('tr');
    let nombre = document.createElement('td');
    let fecha = document.createElement('td');
    let especialista = document.createElement('td');
    let tipo = document.createElement('td');
    let acciones = document.createElement('td');

    let accionesButton = document.createElement('button');
    accionesButton.classList = "btn btn-primary";
    accionesButton.type = "button";
    accionesButton.textContent = "Consultar entrada";
    acciones.append(accionesButton);
    nombre.textContent = "";
    fecha.textContent = entrada.createdAt;
    especialista.textContent = entrada.empleadoNombre;
    tipo.textContent = "";
    tr.append(nombre);
    tr.append(fecha);
    tr.append(especialista);
    tr.append(tipo);
    tr.append(acciones);
    tr.setAttribute("id", historial.id);
    //Añadimos fila a la tabla
    document.querySelector(`#historialTabla`).querySelector('tbody').append(tr);
}

document.addEventListener('DOMContentLoaded',init,false);