import { Form, Input, Button } from "@arco-design/web-react";
import { login } from "../services/api";
const FormItem = Form.Item;

function Login() {
  const [form] = Form.useForm();
  return (
    <Form
      form={form}
      style={{ width: 600 }}
      initialValues={{ name: "admin" }}
      autoComplete="off"
      onValuesChange={(v, vs) => {
        // console.log(v, vs);
      }}
      onSubmit={(v) => {
        console.log(v);
        login(v.username, v.password);
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
