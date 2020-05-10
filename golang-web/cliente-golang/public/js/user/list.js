var SELECTED_USER_ID;

function loadTable(userList){
    if(!userList || userList.length < 1){
        document.querySelector("#alert").textContent = "No hay resultados";
        document.querySelector("#alert").classList.replace("alert-success", "alert-danger");
        document.querySelector("#alert").classList.remove('invisible');
    }
    userList.forEach(user => {
        addRow(user);
    });
}

function addRow(user){
    let tr = document.createElement('tr');
    let id = document.createElement('td');
    let name = document.createElement('td');
    let surname = document.createElement('td');
    let email = document.createElement('td');
    let date = document.createElement('td');
    let actions = document.createElement('td');

    let editButton = document.createElement('button');
    editButton.classList = "btn btn-primary";
    editButton.type = "button";
    editButton.textContent = "Editar usuario";
    editButton.addEventListener("click", editUser, false);
    actions.append(editButton);
    let deleteButton = document.createElement('button');
    deleteButton.setAttribute("data-toggle", "modal");
    deleteButton.setAttribute("data-target", "#borradoModal");
    deleteButton.classList = "btn btn-danger";
    deleteButton.type = "button";
    deleteButton.textContent = "Eliminar usuario";
    deleteButton.addEventListener("click", selectDeleteUser, false);
    actions.append(deleteButton);
    id.textContent = user.id;
    name.textContent = user.nombre;
    surname.textContent = user.apellidos;
    email.textContent = user.email;
    date.textContent = user.createdAt.replace('T',' ').replace('Z','');
    tr.append(id);
    tr.append(name);
    tr.append(surname);
    tr.append(email);
    tr.append(date);
    tr.append(actions);
    tr.setAttribute("id", user.id);
    //Añadimos fila a la tabla
    document.querySelector(`#tablaUsuarios`).querySelector('tbody').append(tr);
}

function editUser(event){
    userId = event.target.closest("tr").getAttribute("id");
    document.location.href="/user/adminG/userList/" + userId + "/edit";
}

function selectDeleteUser(event){
    SELECTED_USER_ID = event.target.closest("tr").getAttribute("id");
}

function deleteUser(event){
    deleteUserREST(SELECTED_USER_ID);
}

function deleteUserREST(user_id){
    const url= "/user/" + user_id + "/delete";
    const request = {
        method: 'DELETE', 
        headers: cabeceras,
    };
    fetch(url,request)
    .then( response => response.json() )
        .then( r => {
            if(!r.Error){
                location.reload();
            }
            else{
                alert(r.Error)
            }
        })
        .catch(err => alert(err));
}

function init () {
    deleteBreadcrumb();
    addLinkBreadcrumb('Usuario', '/user/menu');
    addLinkBreadcrumb('Administrador global', '/user/adminG');
    addLinkBreadcrumb('Usuarios', '/user/adminG/userList');
    loadTable(userList);
    modalBorrar = document.querySelector("#borrarUserIdModal");
    modalBorrar.addEventListener("click", deleteUser, false);
}

document.addEventListener('DOMContentLoaded',init,false);