# 导入必要的库
from flask import Flask, request, jsonify
import ollama

# 创建Flask应用
app = Flask(__name__)

# 定义POST端点
@app.route('/chat', methods=['POST'])
def chat():
    try:
        # 从请求中获取消息
        data = request.json
        msg = data.get('msg', '')

        # 使用ollama处理消息
        response = ollama.chat(model='llama3.1:8b', messages=[
            {
                'role': 'user',
                'content': msg,
            },
        ])

        # 返回响应
        return jsonify({'response': response['message']['content']})
    except ollama._types.ResponseError as e:
        return jsonify({'error': f"An error occurred: {e.text} (Status code: {e.status_code})"}), 500

# 运行Flask应用
if __name__ == '__main__':
    app.run(debug=True)