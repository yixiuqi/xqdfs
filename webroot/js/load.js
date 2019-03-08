var _PageHeight = document.documentElement.clientHeight,
_PageWidth = document.documentElement.clientWidth;
//计算loading框距离顶部和左部的距离（loading框的宽度为215px，高度为61px）
var _LoadingTop = _PageHeight > 61 ? (_PageHeight - 61) / 2 : 0,
_LoadingLeft = _PageWidth > 215 ? (_PageWidth - 215) / 2 : 0;
//加载gif地址
var Loadimagerul="themes/ui/images/loading.gif";
//在页面未加载完毕之前显示的loading Html自定义内容
var _LoadingHtml = '<div id="loadingDiv" style="position:absolute;left:0;width:100%;height:' + _PageHeight + 'px;top:0;background:#f3f8ff;opacity:1;filter:alpha(opacity=80);z-index:10000;"><div style="position: absolute; cursor1: wait; left: ' + _LoadingLeft + 'px; top:' + _LoadingTop + 'px; width:130px;height: 32px; line-height: 32px; padding-left: 27px; background: #fff url('+Loadimagerul+') no-repeat scroll 5px 9px; border: 1px solid #E2E2E2; color: #696969;">loading page...</div></div>';
//呈现loading效果
document.write(_LoadingHtml);
//监听加载状态改变
document.onreadystatechange = completeLoading;
//加载状态为complete时移除loading效果
function completeLoading() {
	if (document.readyState == "complete") {
		var loadingMask = document.getElementById('loadingDiv');
		loadingMask.parentNode.removeChild(loadingMask);
	}
}
