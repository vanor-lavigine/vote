import React, { useEffect, useState } from "react";
import { observer } from "mobx-react-lite";
import UserStore from "../store";
import {
  listCandidates,
  deleteCandidate,
  createCandidate,
} from "../services/api";
import { Table, Button, Modal, Input, Form } from "@arco-design/web-react";

interface Candidate {
  id: number;
  username: string;
  voted: boolean;
}

const CandidateList: React.FC = observer(() => {
  const [candidates, setCandidates] = React.useState<Candidate[]>([]);
  const [showModal, setShowModal] = useState(false);
  const [form] = Form.useForm();

  const fetchCandidates = async () => {
    const candidates :any = await listCandidates();
    //const list = candidates?.data || [];
   // console.log("候选人列表", list);
    setCandidates(candidates);
    //console.log(candidates);
  };

  useEffect(() => {
    fetchCandidates();
  }, []);
  const handleDelete = async (id: number) => {
    await deleteCandidate(id);
    fetchCandidates();
  };

  const handleCreate = async () => {
    form.validate().then(
      async () => {
        const username = form.getFieldValue("Username");
        try {
          await createCandidate(username);
          fetchCandidates();
          setShowModal(false);
        } catch (error) {
          console.log("创建失败", error);
        }
      },
      (err) => {
        console.log(err);
      }
    );
  };

  const columns = [
    {
      title: "ID",
      dataIndex: "id",
      key: "Id",
    },
    {
      title: "用户名",
      dataIndex: "username",
      key: "Username",
    },
    {
      title: "操作",
      key: "action",
      render: (_:any, record:Candidate) => (
        <>
          {UserStore.username === "admin" && (
            <Button onClick={() => handleDelete(record.id)} status="danger">
              Delete
            </Button>
          )}
        </>
      ),
    },
  ];

  return (
    <div>
      <Button type="primary" onClick={() => setShowModal(true)}>创建候选人</Button>
      <p></p>
      <Modal
        visible={showModal}
        onCancel={() => setShowModal(false)}
        title="创建候选人"
        onOk={handleCreate}
        okText="创建"
      >
        <Form form={form}>
          <Form.Item label="用户名" field="Username" required>
            <Input placeholder="请输入用户名" />
          </Form.Item>
        </Form>
      </Modal>
      <Table columns={columns} data={candidates} rowKey="id" />
    </div>
  );
});

export default CandidateList;
