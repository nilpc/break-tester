from flask import Flask, jsonify
import time
import random

app = Flask(__name__)

@app.route('/users', methods=['GET'])
def get_users():
    time.sleep(random.uniform(0.05, 0.15))
    return jsonify({"status": "success", "users": ["Alice", "Bob", "Charlie"]})

@app.route('/checkout', methods=['POST'])
def process_checkout():
    time.sleep(random.uniform(0.1, 0.3))
    return jsonify({"status": "success", "message": "Payment processed successfully!"})

if __name__ == '__main__':
    print("🍳 The Kitchen is open for business on http://localhost:5000")
    app.run(host='0.0.0.0', port=5000)