import React, { useEffect, useState } from "react";
import { observer } from "mobx-react-lite";
import UserStore from "../store";
import {listCandidates, listCandidatesVoteStatus, vote} from "../services/api";
import { Table, Button, Message } from "@arco-design/web-react";

interface Candidate {
  id: number;
  username: string;
  voted: boolean;
}

const Vote: React.FC = observer(() => {
  const [candidates, setCandidates] = React.useState<Candidate[]>([]);
  const [voted, setVoted] = useState<boolean>(false);

  useEffect(() => {
    fetchCandidates();
  }, []);

  const fetchCandidates = async () => {
    const candidates = await listCandidatesVoteStatus();
    const list: any = candidates || [];
    console.log("候选人列表", list);

    // 找到一个投过票的，说明已投过，只能投一次
    const _hasVoted = !!list?.find((c:any) => c.voted);
    if (_hasVoted) {
      Message.info("您已投过票，无法再投");
    }
    setVoted(_hasVoted);
    setCandidates(list);
  };

  const handleVote = async (username: string) => {
    try {
      await vote(username);
      Message.info("投票成功");
      fetchCandidates();
    } catch (e) {
      console.error('投票失败', e);
      Message.info("投票失败");
    }
    // 投票
    console.log("投票", username);
  };

  const columns = [
    {
      title: "ID",
      dataIndex: "Id",
      key: "id",
    },
    {
      title: "用户名",
      dataIndex: "username",
      key: "username",
    },
    {
      title: "投票",
      key: "action",
      render: (record: Candidate) => (
          <Button
              disabled={voted}
              onClick={() => handleVote(record.username)}
          >
            投票
          </Button>
      ),
    },
  ];

  return (
    <div>
      <Table columns={columns} data={candidates} rowKey="id" />
    </div>
  );
});

export default Vote;
