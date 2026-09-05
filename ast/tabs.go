package ast

// TabTitleBlock 返回作为页签标题保存的原段落，段落仍具有独立的块身份。
func (n *Node) TabTitleBlock() *Node {
	if nil == n || NodeTabItem != n.Type {
		return nil
	}
	for child := n.FirstChild; nil != child; child = child.Next {
		if NodeKramdownBlockIAL == child.Type {
			continue
		}
		if NodeParagraph == child.Type && "true" == child.IALAttr("tabs-title") {
			return child
		}
		break
	}
	return nil
}

func (n *Node) IsTabTitleBlock() bool {
	return nil != n.Parent && n.Parent.TabTitleBlock() == n
}
