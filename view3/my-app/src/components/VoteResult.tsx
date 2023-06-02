import React, { useState, useEffect } from "react";
import { observer } from "mobx-react-lite";
import { getVoteResult} from "../services/api";
import { Table } from "@arco-design/web-react";

interface VoteItemResult {
  id: number;
  username: string;
  count: number;
}

const VoteResult: React.FC = observer(() => {
  const [voteList, setVoteList] = React.useState<VoteItemResult[]>([]);

  useEffect(() => {
    fetchVoteList();
  }, []);

  const fetchVoteList = async () => {
    const list:any = await getVoteResult();
    setVoteList(list || []);
  };

  const columns = [
    {
      title: "id",
      dataIndex: "id",
      key: "id",
    },
    {
      title: "用户名",
      dataIndex: "username",
      key: "username",
    },
    {
      title: "票数",
      dataIndex: "count",
      key: "count",
    },
  ];

  return <Table columns={columns} data={voteList} rowKey="id" />;
});

export default VoteResult;
