import { Form, Input, Button } from "@arco-design/web-react";
import { login } from "../services/api";
import {useHistory} from "react-router-dom";
import React, { useEffect, useState, useRef } from 'react';
import { Evaluate, ProofHoHash } from '@idena/vrf-js'


const FormItem = Form.Item;

interface VerificationCodeImageProps {
  code: any; // 你可以将 'any' 替换为实际的类型
}

function VerificationCodeImage({ code }: VerificationCodeImageProps) {
  const canvasRef = useRef<HTMLCanvasElement | null>(null);

  useEffect(() => {
    const canvas = canvasRef.current;
    if (canvas) {
      const ctx = canvas.getContext('2d');

      if (ctx) {
        // 清空画布
        ctx.clearRect(0, 0, canvas.width, canvas.height);

        // 在清空的画布上添加新的验证码
        ctx.font = '20px Arial';
        ctx.fillText(code, 10, 50);
      }
    }
  }, [code]);

  return <canvas ref={canvasRef} />;
}


function Login() {
  const [form] = Form.useForm();
  const history = useHistory();

  const [code, setCode] = useState(''); // 添加一个状态来存储验证码
  const privateKeys = [
    [123, 254, 12, /*...*/ 11],
    [120, 254, 13, /*...*/ 12],
    // ...更多私钥
  ];
  const dataSets = [
    [1, 2, 3, 4, 5],
    [2, 3, 4, 5, 6],
    // ...更多数据集
  ];
  const generateNewCode = () => {
    const privateKeyIndex = Math.floor(Math.random() * privateKeys.length); // 随机挑选一个私钥
    const privateKey = privateKeys[privateKeyIndex];

    const dataIndex = Math.floor(Math.random() * dataSets.length); // 随机挑选一个数据数组
    const data = dataSets[dataIndex];

    const [hash, proof] = Evaluate(privateKey, data);
    const hashHex = hash.map((byte: number) => ('0' + byte.toString(16)).slice(-2)).join('');

    let newCode = hashHex.slice(-6);
    setCode(newCode); // 更新验证码的状态
    console.log(newCode)
  }


  useEffect(() => {
    generateNewCode();
  }, []);



  return (
    <Form
      form={form}
      style={{ width: 600 }}
      initialValues={{ name: "admin" }}
      autoComplete="off"
      onValuesChange={(v, vs) => {
        // console.log(v, vs);
      }}
      onSubmit={async (v) => {
        console.log(v);
        try {
         await login(v.username, v.password);
         history.push('/');
        } catch (e) {
          console.log('login error')
        }



      }}
    >
      <FormItem label="用户名" field="username" rules={[{ required: true }]}>
        <Input placeholder="请输入用户名" />
      </FormItem>
      <FormItem label="密码" field="password" rules={[{ required: true }]}>
        <Input.Password placeholder="请输入密码" />
      </FormItem>
      <FormItem wrapperCol={{ offset: 5 }}>
        <Button type="primary" htmlType="submit" style={{ marginRight: 24 }}>
          登录
        </Button>
      </FormItem>
    </Form>
  );
}

export default Login;
