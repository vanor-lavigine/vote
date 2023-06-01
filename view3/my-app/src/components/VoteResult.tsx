import React, { useState, useEffect } from "react";
import { observer } from "mobx-react-lite";
import { getVoteList } from "../services/api";
import { Table } from "@arco-design/web-react";

interface VoteItemResult {
  Id: number;
  Username: string;
  Votes: number;
}

const VoteResult: React.FC = observer(() => {
  const [voteList, setVoteList] = React.useState<VoteItemResult[]>([]);

  useEffect(() => {
    fetchVoteList();
  }, []);

  const fetchVoteList = async () => {
    const list = await getVoteList();
    setVoteList(list?.data || []);
  };

  const columns = [
    {
      title: "ID",
      dataIndex: "Id",
      key: "Id",
    },
    {
      title: "用户名",
      dataIndex: "Username",
      key: "Username",
    },
    {
      title: "票数",
      dataIndex: "Count",
      key: "Count",
    },
  ];

  return <Table columns={columns} data={voteList} rowKey="id" />;
});

export default VoteResult;
