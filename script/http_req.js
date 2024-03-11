async function sendSynchronously() {
    try {
        const response = await fetch('http://www.google.com/');
        if (!response.ok) {
            throw new Error(`HTTP 错误！状态码: ${response.status}`);
        }
        const data = await response.json();
        // 处理同步返回的结果
        console.log(data);
    } catch (error) {
        // 处理错误
        console.error('请求失败:', error);
    }
}

function run(input){
    sendSynchronously()
}
