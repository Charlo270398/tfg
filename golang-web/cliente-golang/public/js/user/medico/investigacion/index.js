function init () {
    deleteBreadcrumb();
    addLinkBreadcrumb('Usuario', '/user/menu');
    addLinkBreadcrumb('Medico', '/user/doctor');
    addLinkBreadcrumb('Investigación', '');
    console.log(analiticas);
}

document.addEventListener('DOMContentLoaded',init,false);