
function loadTable(clinicaList){
    if(!clinicaList || clinicaList.length < 1){
        document.querySelector("#alert").textContent = "No hay resultados";
        document.querySelector("#alert").classList.replace("alert-success", "alert-danger");
        document.querySelector("#alert").classList.remove('invisible');
        document.querySelector("#tablaClinicas").classList.add('invisible');
    }else{
        clinicaList.forEach(clinica => {
            addRow(clinica);
        });
    }
}

function addRow(clinica){
    let tr = document.createElement('tr');
    let id = document.createElement('td');
    let name = document.createElement('td');
    let medicos = document.createElement('td');
    let enfermeros = document.createElement('td');
    let administradores = document.createElement('td');
    let actions = document.createElement('td');

    let editButton = document.createElement('button');
    editButton.classList = "btn btn-primary";
    editButton.type = "button";
    editButton.textContent = "Editar clinica";
    editButton.addEventListener("click", editClinica, false);
    actions.append(editButton);
    let deleteButton = document.createElement('button');
    deleteButton.classList = "btn btn-danger";
    deleteButton.type = "button";
    deleteButton.textContent = "Eliminar clinica";
    deleteButton.addEventListener("click", deleteClinica, false);
    actions.append(deleteButton);

    id.textContent = clinica.Id;
    name.textContent = clinica.Nombre;
    medicos.textContent = clinica.NumeroMedicos;
    enfermeros.textContent = clinica.NumeroEnfermeros;
    administradores.textContent = clinica.NumeroAdministradores;
    tr.append(id);
    tr.append(name);
    tr.append(enfermeros);
    tr.append(medicos);
    tr.append(administradores);
    tr.append(actions);
    tr.setAttribute("id", clinica.Id);
    //Añadimos fila a la tabla
    document.querySelector(`#tablaClinicas`).querySelector('tbody').append(tr);
}

function editClinica(event){
    clinicaId = event.target.closest("tr").getAttribute("id");
    document.location.href="/clinica/list/" + clinicaId + "/edit";
}

function deleteClinica(event){
    clinicaId = event.target.closest("tr").getAttribute("id");
    //AQUI REST DELETE PERO PRIMERO UN MODAL
}

function init () {
    deleteBreadcrumb();
    addLinkBreadcrumb('Usuario', '/user/menu');
    addLinkBreadcrumb('Administrador global', '/user/adminG');
    addLinkBreadcrumb('Clinicas', '/clinica/list');
    loadTable(clinicaList);
}

document.addEventListener('DOMContentLoaded',init,false);