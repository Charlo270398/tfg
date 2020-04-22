function loadTable(especialidadList){
    if(!especialidadList || especialidadList.length < 1){
        document.querySelector("#alert").textContent = "No hay resultados";
        document.querySelector("#alert").classList.replace("alert-success", "alert-danger");
        document.querySelector("#alert").classList.remove('invisible');
        document.querySelector("#tablaEspecialidades").classList.add('invisible');
    }else{
        especialidadList.forEach(especialidad => {
            addRow(especialidad);
        });
    }
}

function addRow(especialidad){
    let tr = document.createElement('tr');
    let id = document.createElement('td');
    let name = document.createElement('td');
    let actions = document.createElement('td');

    let editButton = document.createElement('button');
    editButton.classList = "btn btn-primary";
    editButton.type = "button";
    editButton.textContent = "Editar especialidad";
    editButton.addEventListener("click", editEspecialidad, false);
    actions.append(editButton);
    let deleteButton = document.createElement('button');
    deleteButton.classList = "btn btn-danger";
    deleteButton.type = "button";
    deleteButton.textContent = "Eliminar especialidad";
    deleteButton.addEventListener("click", deleteEspecialidad, false);
    actions.append(deleteButton);

    id.textContent = especialidad.Id;
    name.textContent = especialidad.Nombre;
    tr.append(id);
    tr.append(name);
    tr.append(actions);
    tr.setAttribute("id", especialidad.Id);
    //Añadimos fila a la tabla
    document.querySelector(`#tablaEspecialidades`).querySelector('tbody').append(tr);
}

function editEspecialidad(event){
    especialidadId = event.target.closest("tr").getAttribute("id");
    document.location.href="/especialidad/list/" + especialidadId + "/edit";
}

function deleteEspecialidad(event){
    clinicaId = event.target.closest("tr").getAttribute("id");
    //AQUI REST DELETE PERO PRIMERO UN MODAL
}

function init () {
    deleteBreadcrumb();
    addLinkBreadcrumb('Usuario', '/user/menu');
    addLinkBreadcrumb('Administrador global', '/user/adminG');
    addLinkBreadcrumb('Especialidades', '/especialidad/list');
    loadTable(especialidadList);
}

document.addEventListener('DOMContentLoaded',init,false);