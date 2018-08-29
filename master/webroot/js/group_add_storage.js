function httpGroupAddStorage(groupId,storageId,storageAddr){
    var str="groupId="+groupId+"&storageId="+storageId+"&storageAddr="+storageAddr;
    var result
    var error="null"
    $.ajax({
        url: "/group/storage/add",  
        type: "POST",
        data: str,
        async:false, 
        dataType:"json", 
        success: function( response ) {
            if(response.result==0){
                result=0
            }else{
                result=response.result
                error=response.info
            }            
        },
        error: function( jqXHR, textStatus, errorThrown ) {
            result=-1
        }
    });

    return [result,error]
}
