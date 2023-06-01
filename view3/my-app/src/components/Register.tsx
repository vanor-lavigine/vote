import { Form, Input, Button } from "@arco-design/web-react";
import { register } from "../services/api";
const FormItem = Form.Item;

function Register() {
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
        register(v.username, v.password);
      }}
    >
      <FormItem label="用户名" field="username" rules={[{ required: true }]}>
        <Input placeholder="请输入用户名" />
      </FormItem>
      <FormItem
        label="密码"
        field="password"
        rules={[{ required: true, minLength: 8 }]}
      >
        <Input.Password placeholder="请输入密码" />
      </FormItem>
      <FormItem wrapperCol={{ offset: 5 }}>
        <Button type="primary" htmlType="submit" style={{ marginRight: 24 }}>
          注册
        </Button>
      </FormItem>
    </Form>
  );
}

export default Register;
