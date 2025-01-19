import { Box, Skeleton, styled } from "@mui/material";

const BookCard = styled(Box)<{ $width: number }>(({ $width }) => ({
  display: "flex",
  flexDirection: "column",
  alignItems: "center",
  width: $width,
  gap: 8,
}));

type BookPreviewSkeletonProps = {
  width?: number;
};

export const BookPreviewSkeleton = ({ width = 150 }: BookPreviewSkeletonProps) => {
  return (
    <BookCard $width={width}>
      <Skeleton variant="rectangular" width={width} height={width * 1.4} sx={{ borderRadius: 2 }} />
      <Skeleton variant="text" width="100%" height={32} />
    </BookCard>
  );
};
